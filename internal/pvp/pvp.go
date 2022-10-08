package pvp

import (
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/logrusx"
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
			logrusx.Log.Info("pvp app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.Info("pvp app exit .....")
		}),
		app.Name("pvp"),
		app.Runners(runners...),
	)
	if err := pvp.Run(); err != nil {
		return err
	}
	return nil
}
