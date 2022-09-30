package pvp

import (
	"log"

	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
)

type PvpApp struct {
}

func NewPvpApp() *PvpApp {
	return &PvpApp{}
}

func (p *PvpApp) Run(cfg config.GameConfig) error {
	runners := make([]app.Runner, 0)
	pvp := app.New(
		app.OnBeginHook(func() {
			log.Println("pvp app start....")
		}),
		app.OnExitHook(func() {
			log.Println("pvp app exit....")
		}),
		app.Name("pvp"),
		app.Runners(runners...),
	)
	if err := pvp.Run(); err != nil {
		return err
	}
	return nil
}
