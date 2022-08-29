package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/junaozun/game_server/internal/battle"
	"github.com/junaozun/game_server/internal/cross"
	"github.com/junaozun/game_server/internal/logic"
	"github.com/junaozun/game_server/internal/pvp"
	"github.com/junaozun/game_server/internal/web"
	"github.com/junaozun/game_server/pkg/app"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
)

// Service 接口
type Service interface {
	Init(pkgConfig.GameConfig) error
	ParseFlag(*flag.FlagSet)
	app.Runner
}

var (
	cfgPath = flag.String("config", "game.yaml", "config file path")
)

func main() {
	cfg := pkgConfig.GameConfig{}
	if err := pkgConfig.LoadConfigFromFile(*cfgPath, &cfg); nil != err {
		panic(err)
	}
	var servers = []Service{logic.NewLogicService(), battle.NewBattleService(), cross.NewCrossService(), pvp.NewPvpService()}
	runners := make([]app.Runner, 0, len(servers))
	for _, srv := range servers {
		srv.ParseFlag(flag.CommandLine)
		err := srv.Init(cfg)
		if err != nil {
			panic(err)
		}
		runners = append(runners, srv)
	}
	runners = append(runners, web.NewHttpServiceWithConfig(cfg.Web, cfg.Common))
	rand.Seed(time.Now().UnixNano())
	app := app.New(
		app.Name("sanguo"),
		app.Version("v1.0"),
		app.Runners(runners...),
	)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
