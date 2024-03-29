package ws

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/arl/statsviz"
	"github.com/gorilla/websocket"
)

type WsServer struct {
	*http.Server
	ServerRouter *Router
	isGateway    bool
}

func NewWsServer(addr string, router *Router, isGateway bool) *WsServer {
	return &WsServer{
		Server:       &http.Server{Addr: addr},
		ServerRouter: router,
		isGateway:    isGateway,
	}
}

func (s *WsServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.wsHandler)
	// 默认情况下statsviz的服务路由地址是在 /debug/statsviz/下
	// http://localhost:8002/debug/statsviz/
	statsviz.Register(mux)
	// statsviz.RegisterDefault()
	s.Handler = mux
	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *WsServer) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

var wsUpgreader = websocket.Upgrader{
	// 允许所有的CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	// 将http协议升级websocket
	wsConn, err := wsUpgreader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("websocket服务连接失败")
	}
	// websocket通道建立之后 不管是客户端还是服务端 都可以收发消息
	// 发消息的时候把消息当做路由来处理 消息是有格式的 先定义消息的格式
	// 客户端发消息的时候 {Name:"account.login"} 收到之后进行解析，认为想要处理登录逻辑
	wsServer := newWsServer(wsConn, s.isGateway)
	wsServer.addRouter(s.ServerRouter)
	wsServer.start()
}
