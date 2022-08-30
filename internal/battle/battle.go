package battle

import (
	"context"
	"log"

	"github.com/junaozun/game_server/pkg/app"
	"github.com/junaozun/game_server/pkg/config"
)

type BattleApp struct {
}

func NewBattleApp() *BattleApp {
	return &BattleApp{}
}

func (b *BattleApp) Run(ctx context.Context, cfg config.GameConfig) error {
	runners := make([]app.Runner, 0)
	cross := app.New(
		app.OnBeginHook(func() {
			log.Println("battle app start....")
		}),
		app.OnExitHook(func() {
			log.Println("battle app exit....")
		}),
		app.Runners(runners...),
		app.Name("battle"),
	)
	if err := cross.Run(); err != nil {
		return err
	}
	return nil
}
