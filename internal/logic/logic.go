package logic

import (
	"flag"

	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/logic/game"
	"github.com/junaozun/game_server/internal/logic/nats_handler"
	"github.com/junaozun/game_server/internal/logic/wsMgr"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/dao"
	"github.com/junaozun/gogopkg/logrusx"
	"github.com/junaozun/gogopkg/natsx"
)

var (
	ServerId string
)

const (
	host = "0.0.0.0:"
)

type LogicApp struct {
	onLineUser *wsMgr.WsMgr
	ServerName string
}

func NewLogicApp() *LogicApp {
	return &LogicApp{
		onLineUser: wsMgr.NewWsMgr(),
		ServerName: common.ServerName_Logic,
	}
}

func (l *LogicApp) Run(cfg config.GameConfig) error {
	flag.CommandLine.StringVar(&ServerId, "server_id", "100001", "logic server id")
	dao, err := dao.NewDao([]interface{}{cfg.Logic.Mysql, cfg.Common.Etcd, cfg.Common.Redis})
	if err != nil {
		panic(err)
	}
	wsServer := ws.NewWsServer(host+cfg.Logic.Port, ws.NewRouter(), false)
	// 初始化游戏玩法
	game.NewGame(component.NewComponent(dao, cfg), wsServer.ServerRouter, l.onLineUser)

	natsxServer := natsx.New(cfg.Common.NATS, l.ServerName)
	// 注册nats
	nats_handler.RegisterHandler(natsxServer, ServerId)

	logic := app.New(
		app.OnBeginHook(func() {
			logrusx.Log.WithFields(logrusx.Fields{
				"addr": wsServer.Addr,
			}).Info("logic app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.WithFields(logrusx.Fields{
				"addr": wsServer.Addr,
			}).Info("logic app exit .....")
		}),
		app.Name(l.ServerName+ServerId),
		app.Runners(wsServer, natsxServer),
	)
	if err := logic.Run(); err != nil {
		return err
	}
	return nil
}
