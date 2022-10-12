package rank

import (
	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/internal/rank/rank_server"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/logrusx"
	"github.com/junaozun/gogopkg/natsx"
)

type RankApp struct {
	ServerName string
}

func NewRankApp() *RankApp {
	return &RankApp{
		ServerName: common.ServerName_Rank,
	}
}

func (c *RankApp) Run(cfg config.GameConfig) error {
	rankServer := rank_server.NewRank()

	natsxServer := natsx.New(cfg.Common.NATS, c.ServerName)
	// 注册nats

	natsxServer.Register(natsxServer.ServerName, &rank_server.RankHandler{
		Rank: rankServer,
	})

	rank := app.New(
		app.OnBeginHook(func() {
			logrusx.Log.Info("rank app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.Info("rank app exit .....")
		}),
		app.Name(c.ServerName),
		app.Runners(natsxServer, rankServer),
	)
	if err := rank.Run(); err != nil {
		return err
	}
	return nil
}
