package chess

import (
	"log"

	"github.com/junaozun/game_server/internal/cross/gvg/nats_handler"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/natsx"
)

type GvgApp struct {
	ServerName string
}

func NewGvgApp() *GvgApp {
	return &GvgApp{
		ServerName: "gvg",
	}
}

func (c *GvgApp) Run(cfg config.GameConfig) error {
	runners := make([]app.Runner, 0)
	natsxServer := natsx.New(cfg.Common.NATS, c.ServerName)
	// 注册nats
	nats_handler.RegisterHandler(natsxServer)
	runners = append(runners, natsxServer)
	gvg := app.New(
		app.OnBeginHook(func() {
			log.Println("gvg app start....")
		}),
		app.OnExitHook(func() {
			log.Println("gvg app exit....")
		}),
		app.Name("gvg"),
		app.Runners(runners...),
	)
	if err := gvg.Run(); err != nil {
		return err
	}
	return nil
}
