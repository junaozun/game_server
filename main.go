package main

import (
	"github.com/junaozun/game_server/config"
	"github.com/junaozun/game_server/net"
)

func main() {
	host := config.File.MustValue("login_server", "host", "0.0.0.0")
	port := config.File.MustValue("login_server", "port", "8003")
	server := net.NewServer(host + ":" + port)
	server.Start()
}
