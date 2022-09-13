package chess

import (
	"log"

	"github.com/junaozun/game_server/internal/cross/gvg/nats_handler"
	"github.com/junaozun/game_server/pkg/app"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
)

type GvgApp struct {
}

func NewGvgApp() *GvgApp {
	return &GvgApp{}
}

func (c *GvgApp) Run(cfg pkgConfig.GameConfig) error {
	runners := make([]app.Runner, 0)
	// 注册nats
	nats_handler.RegisterHandler("gvg", cfg.Common.NATS)
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
