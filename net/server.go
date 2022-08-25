package net

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	*http.Server
	ServerRouter *Router
}

func NewServer(addr string, router *Router) *Server {
	return &Server{
		Server:       &http.Server{Addr: addr},
		ServerRouter: router,
	}
}

func (s *Server) Start(ctx context.Context) error {
	http.HandleFunc("/", s.wsHandler)
	log.Printf("logic server start success,listenAddr:%s", s.Addr)
	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		log.Println("[httpServer] stop violence.........")
		return err
	}
	log.Println("[httpServer] stop elegant.........")
	return nil
}

var wsUpgreader = websocket.Upgrader{
	// 允许所有的CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// 将http协议升级websocket
	wsConn, err := wsUpgreader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("websocket服务连接失败")
	}
	// websocket通道建立之后 不管是客户端还是服务端 都可以收发消息
	// 发消息的时候把消息当做路由来处理 消息是有格式的 先定义消息的格式
	// 客户端发消息的时候 {Name:"account.login"} 收到之后进行解析，认为想要处理登录逻辑
	wsServer := NewWsServer(wsConn)
	wsServer.AddRouter(s.ServerRouter)
	wsServer.Start()
	// 发送握手协议
	wsServer.Handshake()
}
