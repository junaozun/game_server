package chess

import (
	"fmt"
	"log"

	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/cross/chess/nats_handler"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/dao"
	"github.com/junaozun/gogopkg/natsx"
)

type ChessApp struct {
	ServerName string
}

func NewChessApp() *ChessApp {
	return &ChessApp{
		ServerName: "chess",
	}
}

func (c *ChessApp) Run(cfg config.GameConfig) error {
	runners := make([]app.Runner, 0)
	natsxServer := natsx.New(cfg.Common.NATS, c.ServerName)
	// 注册nats
	nats_handler.RegisterHandler(natsxServer)
	dao, err := dao.NewDao([]interface{}{cfg.Cross.Mysql, cfg.Common.Etcd, cfg.Common.Redis})
	if err != nil {
		panic(err)
	}
	component := component.NewComponent(dao, cfg)
	fmt.Println(component)
	runners = append(runners, natsxServer)
	chess := app.New(
		app.OnBeginHook(func() {
			log.Println("chess app start....")
		}),
		app.OnExitHook(func() {
			log.Println("chess app exit....")
		}),
		app.Name(c.ServerName),
		app.Runners(runners...),
	)
	if err := chess.Run(); err != nil {
		return err
	}
	return nil
}
