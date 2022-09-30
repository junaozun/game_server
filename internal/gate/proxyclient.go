package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/junaozun/game_server/pkg/utils"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/mitchellh/mapstructure"
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

func (p *ProxyClient) Connect() error {
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
		return fmt.Errorf("和服务器 %s 握手失败", p.proxyAddr)
	}
	return nil
}

type syncCtx struct {
	ctx     context.Context // goroutine 的上下文，包含goroutine的运行状态、环境、现场等信息
	cancel  context.CancelFunc
	outChan chan *ws.RespBody
}

func NewSyncCtx() *syncCtx {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	return &syncCtx{
		ctx:     ctx,
		cancel:  cancel,
		outChan: make(chan *ws.RespBody),
	}
}

func (s *syncCtx) wait() *ws.RespBody {
	select {
	case msg := <-s.outChan:
		return msg
	case <-s.ctx.Done():
		log.Println("代理服务响应消息超时")
		return nil
	}
}

type ClientConn struct {
	wsConn        *websocket.Conn // logic 的长链接
	isClosed      bool            // 监听当前客户端是否关闭状态
	handshake     bool            // 握手状态
	handshakeChan chan struct{}   // 接受握手成功信息的通道
	property      map[string]interface{}
	propertyMutex sync.RWMutex
	Seq           int64
	onPush        func(conn *ClientConn, body *ws.RespBody)
	onClose       func(conn *ClientConn)
	syncCtxMap    map[int64]*syncCtx
	syncCtxMutex  sync.RWMutex
}

func NewClientConn(conn *websocket.Conn) *ClientConn {
	return &ClientConn{
		wsConn:        conn,
		handshakeChan: make(chan struct{}),
		property:      make(map[string]interface{}),
		syncCtxMap:    make(map[int64]*syncCtx),
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
	return c.wsConn.RemoteAddr().String()
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
	// 等待服务器握手的消息返回
	go c.wsReadLoop()
	return c.waitHandShake()
}

func (c *ClientConn) waitHandShake() bool {
	// 等待握手的成功 等待握手的消息
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case <-c.handshakeChan:
		log.Println("握手成功")
		return true
	case <-ctx.Done():
		log.Println("握手超时")
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
		_, data, err := c.wsConn.ReadMessage()
		if err != nil {
			continue
		}

		// 收到消息后要解析消息 就是json格式
		// 1 data解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错，非法格式,需要json数据：", err)
			continue
		}
		// 2 前端的消息  加密消息 进行解密
		secretKey, ok := c.GetProperty(ws.SecretKey)
		if !ok {
			log.Println("未设置secretKey值")
			continue
		}
		key := secretKey.(string)
		// 客户端传过来的数据是加密的，需要解密
		realData, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			log.Println("数据格式有误，解密失败：", err)
			continue
		}
		data = realData
		// 3.data 转为body
		respBody := &ws.RespBody{}
		err = json.Unmarshal(data, respBody)
		if err != nil {
			log.Println("数据解析失败", err)
			continue
		}
		// 获取到前端传递的数据了
		// 会收到很多消息，可能是握手，心跳，请求信息(account.login)）
		if respBody.Seq == 0 {
			if respBody.Router == ws.HandshakeMsg {
				// 获取服务器秘钥
				handshake := &ws.Handshake{}
				mapstructure.Decode(respBody.Msg, handshake)
				if handshake.Key != "" {
					c.SetProperty(ws.SecretKey, handshake.Key)
				} else {
					log.Println("[ClientConn] handShake key nil")
				}
				c.handshake = true
				c.handshakeChan <- struct{}{}
			} else {
				if c.onPush != nil {
					c.onPush(c, respBody)
				} else {
					log.Println("not set onPush function")
				}
			}
		} else {
			c.syncCtxMutex.RLock()
			ctx, ok := c.syncCtxMap[respBody.Seq]
			c.syncCtxMutex.RUnlock()
			if ok {
				ctx.outChan <- respBody
			} else {
				log.Println("no seq syncCtx find")
			}
		}

	}
	c.Close()
}

func (c *ClientConn) Close() {
	_ = c.wsConn.Close()
}

func (c *ClientConn) write(body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return
	}
	secretKey, ok := c.GetProperty(ws.SecretKey)
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
	err = c.wsConn.WriteMessage(websocket.BinaryMessage, zipData)
	if err != nil {
		log.Println(err)
	}
	log.Println("[wsServer] write2Client success")
}

func (c *ClientConn) SetOnPush(hookFunc func(conn *ClientConn, body *ws.RespBody)) {
	c.onPush = hookFunc
}

func (c *ClientConn) Send(name string, msg interface{}) (*ws.RespBody, error) {
	c.syncCtxMutex.Lock()
	c.Seq++
	return nil, nil
}
