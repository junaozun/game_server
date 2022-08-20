package main

import (
	"flag"

	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/db"
	"github.com/junaozun/game_server/server"
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
	defer db.Close()
	server := server.NewServer(host+cfg.Server.Port, db.DB)
	// 初始化table
	server.InitTable()
	// 初始化路由
	server.InitRouter()
	// 启动服务器
	server.Start()
}
