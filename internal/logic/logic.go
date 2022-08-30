package logic

import (
	"context"
	"flag"
	"log"

	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/logic/game"
	"github.com/junaozun/game_server/internal/logic/wsMgr"
	"github.com/junaozun/game_server/pkg/app"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/dao"
	"github.com/junaozun/game_server/pkg/ws"
)

var (
	ServerId string
)

const (
	host = "0.0.0.0:"
)

type LogicApp struct {
	component  *component.Component
	onLineUser *wsMgr.WsMgr
}

func NewLogicApp() *LogicApp {
	return &LogicApp{
		onLineUser: wsMgr.NewWsMgr(),
	}
}

func (l *LogicApp) Run(ctx context.Context, cfg pkgConfig.GameConfig) error {
	flag.CommandLine.StringVar(&ServerId, "server_id", "", "logic server id")
	dao, err := dao.NewDao([]interface{}{cfg.Logic.Mysql, cfg.Common.Etcd, cfg.Common.Cache})
	if err != nil {
		return err
	}
	l.component = component.NewComponent(dao)
	wsServer := ws.NewWsServer(host+cfg.Logic.Port, ws.NewRouter())
	// 初始化游戏玩法
	game.NewGame(l.component, wsServer.ServerRouter, l.onLineUser)

	logic := app.New(
		app.OnBeginHook(func() {
			log.Printf("logic app start,addr:%s ....", wsServer.Addr)
		}),
		app.OnExitHook(func() {
			log.Printf("logic app exit,addr:%s ....", wsServer.Addr)
		}),
		app.Name("logic"),
		app.Runners(wsServer),
	)
	if err := logic.Run(); err != nil {
		return err
	}
	return nil
}
