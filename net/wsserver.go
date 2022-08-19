package net

import (
	"sync"

	"github.com/gorilla/websocket"
)

type wsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgResp // 回复给客户端的信息
	seq          int64
	property     map[string]interface{} //
	propertyLock sync.RWMutex
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgResp, 1000),
		property: make(map[string]interface{}),
	}
}

func (w *wsServer) AddRouter(router *Router) {
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

func (w *wsServer) Push(name string, data interface{}) {
	resp := &WsMsgResp{
		Body: &RespBody{
			Seq:  0,
			Name: name,
			Code: 0,
			Msg:  data,
		},
	}
	w.outChan <- resp
}
