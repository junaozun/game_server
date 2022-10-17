package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/junaozun/game_server/internal/gate"
	"github.com/junaozun/game_server/internal/rank"
	"github.com/junaozun/game_server/internal/web"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/logrusx"

	"github.com/junaozun/game_server/internal/battle"
	"github.com/junaozun/game_server/internal/cross/chess"
	"github.com/junaozun/game_server/internal/logic"
	"github.com/junaozun/game_server/internal/pvp"
)

var (
	cfgPath = flag.String("config", "game.yaml", "config file path")
)

func main() {
	go func() {
		logrusx.Log.Info("pprof start.....")
		fmt.Println(http.ListenAndServe(":9876", nil))
	}()
	cfg := config.GameConfig{}
	if err := config.LoadConfigFromFile(*cfgPath, &cfg); nil != err {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	// 将逻辑服、战斗服、跨服、pvp服、web服,gateway,排行榜服全都启动起来
	go logic.NewLogicApp().Run(cfg)
	go battle.NewBattleApp().Run(cfg)
	go chess.NewChessApp().Run(cfg)
	go pvp.NewPvpApp().Run(cfg)
	go web.NewWebApp(cfg).Run()
	go gate.NewGateWay().Run(cfg)
	go rank.NewRankApp().Run(cfg)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	// 主循环
QUIT:
	for {
		select {
		case sig := <-sigs:
			time.Sleep(3 * time.Second)
			log.Printf("Signal: %s", sig.String())
			break QUIT
		case <-ticker.C:
			logrusx.Log.WithFields(logrusx.Fields{
				"goroutine count": runtime.NumGoroutine(),
			}).Info("协程数量")
		}
	}
	logrusx.Log.Info("[main] quiting......")
}
