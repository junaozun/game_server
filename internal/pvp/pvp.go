package pvp

import (
	"context"
	"log"

	"github.com/junaozun/game_server/pkg/app"
	"github.com/junaozun/game_server/pkg/config"
)

type PvpApp struct {
}

func NewPvpApp() *PvpApp {
	return &PvpApp{}
}

func (p *PvpApp) Run(ctx context.Context, cfg config.GameConfig) error {
	runners := make([]app.Runner, 0)
	pvp := app.New(
		app.OnBeginHook(func() {
			log.Println("pvp app start....")
		}),
		app.OnExitHook(func() {
			log.Println("pvp app exit....")
		}),
		app.Runners(runners...),
	)
	if err := pvp.Run(); err != nil {
		return err
	}
	return nil
}
