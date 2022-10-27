package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/pkg/utils"
	"github.com/junaozun/gogopkg/logrusx"
)

var cid int64

const SecretKey = "secretKey"

type wsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgResp // 回复给客户端的信息
	seq          int64
	property     map[string]interface{} // SecretKey:secretKey 、 Cid:cid
	propertyLock sync.RWMutex
	closeWrite   chan struct{}
	isGateway    bool
}

func newWsServer(wsConn *websocket.Conn, isGateway bool) *wsServer {
	s := &wsServer{
		wsConn:     wsConn,
		outChan:    make(chan *WsMsgResp, 1000),
		property:   make(map[string]interface{}),
		closeWrite: make(chan struct{}),
		isGateway:  isGateway,
	}
	if isGateway {
		atomic.AddInt64(&cid, 1)
		// 将自己的cid号设置上
		s.SetProperty("cid", cid)
		// 发送握手协议
		s.handshake()
	} else {
		s.sayHai()
	}
	return s
}

func (w *wsServer) addRouter(router *Router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

func (w *wsServer) GetProperty(key string) (interface{}, bool) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	v, ok := w.property[key]
	return v, ok
}

func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}

func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}

func (w *wsServer) Push(router string, msg interface{}) {
	resp := &WsMsgResp{
		Body: &RespBody{
			Seq:    0,
			Router: router,
			Code:   0,
			Msg:    msg,
		},
	}
	w.outChan <- resp
}

func (w *wsServer) Close() {
	w.closeWrite <- struct{}{}
	_ = w.wsConn.Close()
}

func (w *wsServer) start() {
	// 启动读写数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case wsResp := <-w.outChan:
			w.write2Client(wsResp.Body)
		case <-w.closeWrite:
			return
		}
	}
}

func (w *wsServer) readMsgLoop() {
	defer func() {
		if err := recover(); err != nil {
			w.Close()
			fmt.Printf("[panic] err: %v\nstack: %s\n", err, getCurrentGoroutineStack())
		}
	}()

	// 先读到客户端/gateway发送过来的数据，然后进行处理，再发送回给客户端/gateway消息
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			logrusx.Log.WithFields(logrusx.Fields{"err": err, "isGateway": w.isGateway}).Error("[wsServer] readMsgLoop 收消息出现错误")
			break
		}
		// 收到消息后要解析消息 就是json格式
		// 1 data解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			logrusx.Log.WithFields(logrusx.Fields{"err": err, "isGateway": w.isGateway}).Error("[wsServer] readMsgLoop 解压数据出错，非法格式,需要json数据")
			continue
		}
		// 2 前端的消息  加密消息 进行解密
		// gateway接受client的消息，需要解密;login或logic server接受gateway的消息，不需要解密
		if w.isGateway {
			secretKey, ok := w.GetProperty(SecretKey)
			if !ok {
				logrusx.Log.WithFields(logrusx.Fields{"err": err, "isGateway": w.isGateway}).Error("[wsServer] readMsgLoop 未设置secretKey值")
				continue
			}
			key := secretKey.(string)
			// 客户端传过来的数据是加密的，需要解密
			realData, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				logrusx.Log.WithFields(logrusx.Fields{"err": err, "isGateway": w.isGateway}).Error("[wsServer] readMsgLoop 数据格式有误，解密失败")
				// 出错后发起握手
				w.handshake()
				continue
			}
			data = realData
		}

		// 3.data 转为body
		reqBody := &ReqBody{}
		err = json.Unmarshal(data, reqBody)
		if err != nil {
			log.Println("数据解析失败", err)
			continue
		}

		// 获取到前端传递的数据了，拿上这些数据 去具体的业务进行处理
		wsReq := &WsMsgReq{
			Body: reqBody,
			Conn: w,
		}
		wsRsp := &WsMsgResp{
			Body: &RespBody{
				Seq:    reqBody.Seq,
				Router: reqBody.Router,
			},
		}
		//  交给router处理业务
		w.router.Run(wsReq, wsRsp)
		// 将结果返回给客户端
		w.outChan <- wsRsp
	}
	w.Close()
}

func (w *wsServer) write2Client(resp interface{}) {
	data, err := json.Marshal(resp)
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{"err": err, "isGateway": w.isGateway}).Error("[wsServer] write2Client json marshal err")
		return
	}
	logrusx.Log.WithFields(logrusx.Fields{"data": string(data), "isGateway": w.isGateway}).Info("[wsServer] write2Client replay data")
	// gateway服务器写回客户端的数据需要加密，logic或login 服务器写回gateway的数据只需要压缩就行了
	if w.isGateway {
		secretKey, ok := w.GetProperty(SecretKey)
		if !ok {
			logrusx.Log.WithFields(logrusx.Fields{}).Error("[wsServer] write2Client gateway 未设置secretKey值")
			return
		}
		key := secretKey.(string)
		// 对数据加密
		data, err = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			logrusx.Log.WithFields(logrusx.Fields{"err": err}).Error("[wsServer] write2Client gateway 数据加密错误")
			return
		}
	}

	// 再对数据进行压缩
	zipData, err := utils.Zip(data)
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{"err": err, "isGateway": w.isGateway}).Error("[wsServer] write2Client zip err")
		return
	}
	err = w.wsConn.WriteMessage(websocket.BinaryMessage, zipData)
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{"err": err, "isGateway": w.isGateway}).Error("[wsServer] write2Client writeConn err")
	} else {
		logrusx.Log.WithFields(logrusx.Fields{"isGateway": w.isGateway}).Info("[wsServer] write2Client success")
	}
}

// Handshake 握手协议
// 断开连接重新连接需要进行握手后才能发送消息
// 当游戏客户端 发送请求前先进性一次握手协议、
// 后端会发送对应的加密key给客户端
// 客户端在发送数据的时候，就会使用此key进行加密处理
// 握手协议用于client与gateway第一次连接，gateway生成SecretKey，并将SecretKey发送给客户端，并自己存下
func (w *wsServer) handshake() {
	secretKey := utils.RandSeq(16)
	key, ok := w.GetProperty(SecretKey)
	if ok {
		secretKey = key.(string)
	}
	handshake := &Handshake{Key: secretKey}
	body := &RespBody{Router: common.HandshakeMsg, Msg: handshake}

	data, err := json.Marshal(body)
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{"err": err.Error()}).Error("[wsServer] Handshake json.Marshal err")
		return
	}
	// 数据压缩
	zipData, err := utils.Zip(data)
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{"err": err.Error()}).Error("[wsServer] Handshake zip data err")
		return
	}
	w.wsConn.WriteMessage(websocket.BinaryMessage, zipData)
	// 服务器把key设置上
	w.SetProperty(SecretKey, secretKey)
	logrusx.Log.WithFields(logrusx.Fields{}).Info("[wsServer] Handshake success")
}

func (w *wsServer) sayHai() {
	body := &RespBody{Router: common.SayHaiMsg, Msg: nil}
	data, err := json.Marshal(body)
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{"err": err.Error()}).Error("[wsServer] sayHai json.Marshal err")
		return
	}
	// 数据压缩
	zipData, err := utils.Zip(data)
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{"err": err.Error()}).Error("[wsServer] sayHai zip data err")
		return
	}
	w.wsConn.WriteMessage(websocket.BinaryMessage, zipData)
	logrusx.Log.WithFields(logrusx.Fields{}).Info("[wsServer] sayHai success")
}

// getCurrentGoroutineStack 获取当前Goroutine的调用栈，便于排查panic异常
func getCurrentGoroutineStack() string {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	return string(buf)
}
