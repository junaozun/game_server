package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/junaozun/game_server/pkg/utils"
)

const SecretKey = "secretKey"

type wsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgResp // 回复给客户端的信息
	seq          int64
	property     map[string]interface{} //
	propertyLock sync.RWMutex
	closeWrite   chan struct{}
}

func newWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:     wsConn,
		outChan:    make(chan *WsMsgResp, 1000),
		property:   make(map[string]interface{}),
		closeWrite: make(chan struct{}),
	}
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

func (w *wsServer) Push(router string, data interface{}) {
	resp := &WsMsgResp{
		Body: &RespBody{
			Seq:    0,
			Router: router,
			Code:   0,
			Msg:    data,
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
			w.write2Client(wsResp)
		case <-w.closeWrite:
			return
		}
	}
}

func (w *wsServer) readMsgLoop() {
	defer func() {
		if err := recover(); err != nil {
			w.Close()
			log.Println(err)
		}
	}()

	// 先读到客户端发送过来的数据，然后进行处理，再发送回给客户端消息
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误", err)
			break
		}
		// 收到消息后要解析消息 就是json格式
		// 1 data解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错，非法格式,需要json数据：", err)
			continue
		}
		// 2 前端的消息  加密消息 进行解密
		secretKey, ok := w.GetProperty(SecretKey)
		if !ok {
			log.Println("未设置secretKey值")
			continue
		}
		key := secretKey.(string)
		// 客户端传过来的数据是加密的，需要解密
		realData, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			log.Println("数据格式有误，解密失败：", err)
			// 出错后发起握手
			w.handshake()
			continue
		}
		data = realData
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

func (w *wsServer) write2Client(resp *WsMsgResp) {
	data, err := json.Marshal(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	secretKey, ok := w.GetProperty(SecretKey)
	if !ok {
		log.Println("未设置secretKey值")
		return
	}
	key := secretKey.(string)
	// 对数据加密
	encryptData, err := utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	if err != nil {
		log.Println(err)
		return
	}
	// 再对数据进行压缩
	zipData, err := utils.Zip(encryptData)
	if err != nil {
		log.Println(err)
		return
	}
	err = w.wsConn.WriteMessage(websocket.BinaryMessage, zipData)
	if err != nil {
		log.Println(err)
	}
	log.Println("[wsServer] write2Client success")
}

// Handshake 握手协议
// 断开连接重新连接需要进行握手后才能发送消息
// 当游戏客户端 发送请求前先进性一次握手协议、
// 后端会发送对应的加密key给客户端
// 客户端在发送数据的时候，就会使用此key进行加密处理
func (w *wsServer) handshake() {
	secretKey := utils.RandSeq(16)
	key, ok := w.GetProperty(SecretKey)
	if ok {
		secretKey = key.(string)
	}
	handshake := &Handshake{Key: secretKey}
	body := &RespBody{Router: HandshakeMsg, Msg: handshake}

	data, err := json.Marshal(body)
	if err != nil {
		log.Printf("[wsServer] Handshake json.Marshal err:%s", err.Error())
		return
	}
	// 数据压缩
	zipData, err := utils.Zip(data)
	if err != nil {
		log.Printf("[wsServer] Handshake zip data err:%s", err.Error())
		return
	}
	w.wsConn.WriteMessage(websocket.BinaryMessage, zipData)
	// 服务器把key设置上
	w.SetProperty(SecretKey, secretKey)
}
