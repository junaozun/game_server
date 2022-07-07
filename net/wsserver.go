package net

import (
	"sync"

	"github.com/gorilla/websocket"
)

type wsServer struct {
	wsConn       *websocket.Conn
	router       *router
	outChan      chan *WsMsgResp // 回复给客户端的信息
	seq          int64
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgResp, 1000),
		property: make(map[string]interface{}),
	}
}
