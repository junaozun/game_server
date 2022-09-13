package chess

import (
	"fmt"
	"log"

	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/cross/chess/nats_handler"
	"github.com/junaozun/game_server/pkg/app"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/dao"
)

type ChessApp struct {
}

func NewChessApp() *ChessApp {
	return &ChessApp{}
}

func (c *ChessApp) Run(cfg pkgConfig.GameConfig) error {
	runners := make([]app.Runner, 0)
	// 注册nats
	nats_handler.RegisterHandler("chess", cfg.Common.NATS)
	dao, err := dao.NewDao([]interface{}{cfg.Cross.Mysql, cfg.Common.Etcd, cfg.Common.Cache})
	if err != nil {
		panic(err)
	}
	component := component.NewComponent(dao, cfg)
	fmt.Println(component)
	chess := app.New(
		app.OnBeginHook(func() {
			log.Println("chess app start....")
		}),
		app.OnExitHook(func() {
			log.Println("chess app exit....")
		}),
		app.Name("chess"),
		app.Runners(runners...),
	)
	if err := chess.Run(); err != nil {
		return err
	}
	return nil
}
