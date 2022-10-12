package chess

import (
	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/internal/cross/chess/nats_handler"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/logrusx"
	"github.com/junaozun/gogopkg/natsx"
)

type ChessApp struct {
	ServerName string
}

func NewChessApp() *ChessApp {
	return &ChessApp{
		ServerName: common.ServerName_Chess,
	}
}

func (c *ChessApp) Run(cfg config.GameConfig) error {
	runners := make([]app.Runner, 0)
	natsxServer := natsx.New(cfg.Common.NATS, c.ServerName)
	// 注册nats
	nats_handler.RegisterHandler(natsxServer)
	// dao, err := dao.NewDao([]interface{}{cfg.Cross.Mysql, cfg.Common.Etcd, cfg.Common.Redis})
	// if err != nil {
	// 	panic(err)
	// }
	// component := component.NewComponent(dao, cfg)
	runners = append(runners, natsxServer)
	chess := app.New(
		app.OnBeginHook(func() {
			logrusx.Log.Info("chess app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.Info("chess app exit .....")
		}),
		app.Name(c.ServerName),
		app.Runners(runners...),
	)
	if err := chess.Run(); err != nil {
		return err
	}
	return nil
}
