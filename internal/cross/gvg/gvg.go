package chess

import (
	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/internal/cross/gvg/nats_handler"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/logrusx"
	"github.com/junaozun/gogopkg/natsx"
)

type GvgApp struct {
	ServerName string
}

func NewGvgApp() *GvgApp {
	return &GvgApp{
		ServerName: common.ServerName_Gvg,
	}
}

func (g *GvgApp) Run(cfg config.GameConfig) error {
	runners := make([]app.Runner, 0)
	natsxServer := natsx.New(cfg.Common.NATS, g.ServerName)
	// 注册nats
	nats_handler.RegisterHandler(natsxServer)
	runners = append(runners, natsxServer)
	gvg := app.New(
		app.OnBeginHook(func() {
			logrusx.Log.Info("gvg app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.Info("gvg app exit .....")
		}),
		app.Name(g.ServerName),
		app.Runners(runners...),
	)
	if err := gvg.Run(); err != nil {
		return err
	}
	return nil
}
