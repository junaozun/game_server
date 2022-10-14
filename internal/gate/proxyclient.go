package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/pkg/utils"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/game_server/ret"
	"github.com/junaozun/gogopkg/logrusx"
)

type ProxyClient struct {
	proxyAddr string
	conn      *ClientConn
}

func NewProxyClient(proxy string) *ProxyClient {
	return &ProxyClient{
		proxyAddr: proxy,
	}
}

func (p *ProxyClient) ConnectServer() error {
	// 去连接对应的websocket服务端（可能是login server，也可能是logic server）
	var dialer = websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
		ReadBufferSize:   1023,
		WriteBufferSize:  1024,
		Subprotocols:     []string{"p1", "p2"},
	}
	ws, _, err := dialer.Dial(p.proxyAddr, nil)
	if err != nil {
		return err
	}
	p.conn = NewClientConn(ws)
	if !p.conn.Start() {
		logrusx.Log.WithFields(logrusx.Fields{"server": p.proxyAddr}).Error("gateway 与 服务器握手失败")
		return fmt.Errorf("和服务器 %s 握手失败", p.proxyAddr)
	}
	return nil
}

func (p *ProxyClient) SetProperty(key string, data interface{}) {
	if p.conn.isClosed {
		return
	}
	p.conn.SetProperty(key, data)
}

func (p *ProxyClient) OnPushClient(push func(conn *ClientConn, body *ws.RespBody)) {
	if p.conn.isClosed {
		return
	}
	p.conn.SetOnPush(push)
}

func (p *ProxyClient) Send(router string, msg interface{}) (*ws.RespBody, error) {
	if p.conn.isClosed {
		return nil, fmt.Errorf("conn closed")
	}
	return p.conn.Send(router, msg)
}

// 请求对应的返回
type seqReqRespSync struct {
	ctx     context.Context // goroutine 的上下文，包含goroutine的运行状态、环境、现场等信息
	cancel  context.CancelFunc
	outChan chan *ws.RespBody // 接受logic，login server发来数据的通道
}

func NewSeqReqRespSync() *seqReqRespSync {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	return &seqReqRespSync{
		ctx:     ctx,
		cancel:  cancel,
		outChan: make(chan *ws.RespBody),
	}
}

func (s *seqReqRespSync) wait() *ws.RespBody {
	select {
	case msg := <-s.outChan:
		return msg
	case <-s.ctx.Done():
		logrusx.Log.WithFields(logrusx.Fields{}).Error("代理服务器响应超时")
		return nil
	}
}

// ClientConn 某个用户的客户端连接数据
type ClientConn struct {
	logicWsConn   *websocket.Conn        // gateway与logic的长链接
	isClosed      bool                   // 监听当前客户端是否关闭状态
	sayHai        bool                   // 握手状态
	sayHaiChan    chan struct{}          // 与logic连接成功的信息通道
	property      map[string]interface{} // cid:cid ; proxyAddr:proxyAddr ; clientConn:clientConn
	propertyMutex sync.RWMutex
	Seq           int64
	onPushClient  func(conn *ClientConn, body *ws.RespBody)
	onClose       func(conn *ClientConn)
	seqToResp     map[int64]*seqReqRespSync // seq-> 该序号请求对应的返回
	syncCtxMutex  sync.RWMutex
}

func NewClientConn(conn *websocket.Conn) *ClientConn {
	return &ClientConn{
		logicWsConn: conn,
		sayHaiChan:  make(chan struct{}),
		property:    make(map[string]interface{}),
		seqToResp:   make(map[int64]*seqReqRespSync),
	}
}

func (c *ClientConn) SetProperty(key string, value interface{}) {
	c.propertyMutex.Lock()
	defer c.propertyMutex.Unlock()
	c.property[key] = value
}

func (c *ClientConn) GetProperty(key string) (interface{}, bool) {
	c.propertyMutex.RLock()
	defer c.propertyMutex.RUnlock()
	v, ok := c.property[key]
	return v, ok
}

func (c *ClientConn) RemoveProperty(key string) {
	c.propertyMutex.Lock()
	defer c.propertyMutex.Unlock()
	delete(c.property, key)
}

func (c *ClientConn) Addr() string {
	return c.logicWsConn.RemoteAddr().String()
}

func (c *ClientConn) Push(router string, data interface{}) {
	resp := &ws.WsMsgResp{
		Body: &ws.RespBody{
			Seq:    0,
			Router: router,
			Code:   0,
			Msg:    data,
		},
	}
	c.write(resp.Body)
}

func (c *ClientConn) Start() bool {
	// 做的事情，就是一直不停的接受消息
	// 等待服务器sayHai的消息返回
	go c.wsReadLoop()
	return c.waitSayHai()
}

func (c *ClientConn) waitSayHai() bool {
	// 等待logic服务器回复sayHai消息
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case <-c.sayHaiChan:
		logrusx.Log.WithFields(logrusx.Fields{}).Info("[gateway] ClientConn waitSayHai form logic server success")
		return true
	case err := <-ctx.Done():
		logrusx.Log.WithFields(logrusx.Fields{"err": err}).Info("[gateway] ClientConn waitSayHai form logic server timeout")
		return false
	}
}

// gateway链接logic，循环读取logic发来的数据
func (c *ClientConn) wsReadLoop() {
	defer func() {
		if err := recover(); err != nil {
			c.Close()
			log.Println(err)
		}
	}()
	for {
		_, data, err := c.logicWsConn.ReadMessage()
		if err != nil {
			continue
		}
		// 收到消息后要解析消息 就是json格式
		// 1 data解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			logrusx.Log.WithFields(logrusx.Fields{"err": err}).Error("[gateway] ClientConn wsReadLoop 解压数据出错，非法格式,需要json数据")
			continue
		}

		// 3.data 转为body
		respBody := &ws.RespBody{}
		err = json.Unmarshal(data, respBody)
		if err != nil {
			logrusx.Log.WithFields(logrusx.Fields{"err": err}).Error("[gateway] ClientConn wsReadLoop unmarshal fail")
			continue
		}

		// respBody 为logic发来的数据
		if respBody.Seq == 0 {
			if respBody.Router == common.SayHaiMsg {
				c.sayHai = true
				c.sayHaiChan <- struct{}{}
			} else { // seq == 0却不是sayHai数据，将logic传来的数据原封不动传递给客户端
				if c.onPushClient != nil {
					c.onPushClient(c, respBody)
				} else {
					logrusx.Log.WithFields(logrusx.Fields{}).Error("[gateway] ClientConn wsReadLoop not set onPushClient function ")
				}
			}
		} else {
			c.syncCtxMutex.RLock()
			reqToRespChan, ok := c.seqToResp[respBody.Seq]
			c.syncCtxMutex.RUnlock()
			if ok {
				reqToRespChan.outChan <- respBody
			} else {
				logrusx.Log.WithFields(logrusx.Fields{}).Error("[gateway] ClientConn wsReadLoop no seq seqReqRespSync find ")
			}
		}

	}
	c.Close()
}

func (c *ClientConn) Close() {
	_ = c.logicWsConn.Close()
}

func (c *ClientConn) write(body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return err
	}
	// 再对数据进行压缩
	zipData, err := utils.Zip(data)
	if err != nil {
		log.Println(err)
		return err
	}
	err = c.logicWsConn.WriteMessage(websocket.BinaryMessage, zipData)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("[wsServer] write2Client success")
	return nil
}

func (c *ClientConn) SetOnPush(push func(conn *ClientConn, body *ws.RespBody)) {
	c.onPushClient = push
}

// Send gateway将请求发送给login,logic服务器,等待返回
func (c *ClientConn) Send(router string, msg interface{}) (*ws.RespBody, error) {
	c.syncCtxMutex.Lock()
	c.Seq++
	seq := c.Seq
	sc := NewSeqReqRespSync()
	c.seqToResp[seq] = sc
	c.syncCtxMutex.Unlock()

	req := &ws.ReqBody{
		Seq:    seq,
		Router: router,
		Msg:    msg,
	}

	rsp := &ws.RespBody{
		Seq:    seq,
		Router: router,
		Code:   ret.OK.Code,
	}
	// 将数据写入logic,login 服务器
	err := c.write(req)
	if err != nil {
		sc.cancel()
	} else {
		// 然后等待服务器的返回数据
		r := sc.wait()
		if r == nil {
			rsp.Code = ret.Err_ProxyConnect.Code
		} else {
			rsp = r
		}
	}

	// 该请求处理完成，将请求删除
	c.syncCtxMutex.Lock()
	delete(c.seqToResp, seq)
	c.syncCtxMutex.Unlock()
	return rsp, nil
}
