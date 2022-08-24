package main

import (
	"context"
	"flag"
	"math/rand"
	"time"

	"github.com/junaozun/game_server/internal/battle"
	"github.com/junaozun/game_server/internal/cross"
	"github.com/junaozun/game_server/internal/logic"
	"github.com/junaozun/game_server/internal/pvp"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"golang.org/x/sync/errgroup"
)

var (
	cfgPath = flag.String("config", "game.yaml", "config file path")
)

// Service 接口
type Service interface {
	Init(cfg pkgConfig.GameConfig) error
	ParseFlag(set *flag.FlagSet)
	Run()
}

func main() {

	cfg := pkgConfig.GameConfig{}
	if err := pkgConfig.LoadConfigFromFile(*cfgPath, &cfg); nil != err {
		panic(err)
	}

	var srvs []Service
	srvs = append(srvs,
		logic.NewLogicService(),
		battle.NewBattleService(),
		cross.NewCrossService(),
		pvp.NewPvpService(),
	)

	rand.Seed(time.Now().UnixNano())
	eg, _ := errgroup.WithContext(context.Background())
	for _, v := range srvs {
		srv := v
		eg.Go(func() error {
			err := srv.Init(cfg)
			if err != nil {
				return err
			}
			srv.ParseFlag(flag.CommandLine)
			srv.Run()
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		panic(err)
	}
}
