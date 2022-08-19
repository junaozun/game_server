package main

import (
	"github.com/junaozun/game_server/config"
	"github.com/junaozun/game_server/net"
	"github.com/junaozun/game_server/server"
)

func main() {
	host := config.File.MustValue("login_server", "host", "0.0.0.0")
	port := config.File.MustValue("login_server", "port", "8003")
	addr := host + ":" + port
	s := net.NewServer(addr)
	// 初始化路由
	server.InitRouter(s.ServerRouter)
	s.Start()
}
