package battle

import (
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/logrusx"
)

type BattleApp struct {
}

func NewBattleApp() *BattleApp {
	return &BattleApp{}
}

func (b *BattleApp) Run(cfg config.GameConfig) error {
	runners := make([]app.Runner, 0)
	cross := app.New(
		app.OnBeginHook(func() {
			logrusx.Log.Info("battle app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.Info("battle app exit .....")
		}),
		app.Runners(runners...),
		app.Name("battle"),
	)
	if err := cross.Run(); err != nil {
		return err
	}
	return nil
}
