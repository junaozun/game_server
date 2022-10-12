package gate

import (
	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/logrusx"
)

const (
	host = "127.0.0.1:"
)

type GateApp struct {
	ServerName string
	Router     *ws.Router
	Handler    *Handler
}

func NewGateWay() *GateApp {
	return &GateApp{
		ServerName: common.ServerName_Gate,
		Router:     ws.NewRouter(),
		Handler:    NewHandler(),
	}
}

func (g *GateApp) Run(cfg config.GameConfig) error {
	g.Handler.SetLoginProxy(cfg.GateWay.LoginProxy)
	g.Handler.SetLogicProxy(cfg.GateWay.LogicProxy)
	g.InitRouter()
	wsServer := ws.NewWsServer(host+cfg.GateWay.Port, g.Router)
	gate := app.New(
		app.OnBeginHook(func() {
			logrusx.Log.WithFields(logrusx.Fields{
				"addr": wsServer.Addr,
			}).Info("gate app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.WithFields(logrusx.Fields{
				"addr": wsServer.Addr,
			}).Info("gate app exit .....")
		}),
		app.Name(g.ServerName),
		app.Runners(wsServer),
	)
	if err := gate.Run(); err != nil {
		return err
	}
	return nil
}
