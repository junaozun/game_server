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
	"github.com/junaozun/game_server/internal/cross"
	"github.com/junaozun/game_server/internal/logic"
	"github.com/junaozun/game_server/internal/pvp"
	"github.com/junaozun/game_server/internal/web"
	"github.com/junaozun/game_server/pkg/app"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
)

var (
	cfgPath = flag.String("config", "game.yaml", "config file path")
)

func main() {
	go func() {
		for {
			time.Sleep(time.Second * 3)
			fmt.Printf("协程数量%d", runtime.NumGoroutine())
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
	var apps = []app.IApp{logic.NewLogicApp(), battle.NewBattleApp(), cross.NewCrossApp(), pvp.NewPvpApp(), web.NewWebApp(cfg)}
	appMgr := app.NewAppMgr(apps...)
	appMgr.Runs(cfg)
}
