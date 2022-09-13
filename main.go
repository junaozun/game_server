package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/junaozun/game_server/internal/battle"
	"github.com/junaozun/game_server/internal/cross/chess"
	"github.com/junaozun/game_server/internal/logic"
	"github.com/junaozun/game_server/internal/pvp"
	"github.com/junaozun/game_server/internal/web"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
)

var (
	cfgPath = flag.String("config", "game.yaml", "config file path")
)

func main() {
	go func() {
		for {
			time.Sleep(time.Second * 3)
			fmt.Printf("协程数量%d\n", runtime.NumGoroutine())
		}
	}()
	go func() {
		fmt.Println("pprof start...")
		fmt.Println(http.ListenAndServe(":9876", nil))
	}()
	cfg := pkgConfig.GameConfig{}
	if err := pkgConfig.LoadConfigFromFile(*cfgPath, &cfg); nil != err {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	// 将逻辑服、战斗服、跨服、pvp服、web服全都启动起来
	go logic.NewLogicApp().Run(cfg)
	go battle.NewBattleApp().Run(cfg)
	go chess.NewChessApp().Run(cfg)
	go pvp.NewPvpApp().Run(cfg)
	go web.NewWebApp(cfg).Run()
	for {

	}
}
