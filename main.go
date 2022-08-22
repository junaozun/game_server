package main

import (
	"flag"

	"github.com/junaozun/game_server/internal/game"
	"github.com/junaozun/game_server/net"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/db"
)

var (
	cfgPath = flag.String("config", "game.yaml", "config file path")
)

const host = "0.0.0.0:"

func main() {
	cfg := pkgConfig.Config{}
	if err := pkgConfig.LoadConfigFromFile(*cfgPath, &cfg); nil != err {
		panic(err)
	}
	db, err := db.NewDao(cfg.DB)
	if err != nil {
		panic(err)
	}
	server := net.NewServer(host+cfg.Server.Port, db.DB)
	game.NewGame(server)
	// 启动服务器
	server.Start()
}
