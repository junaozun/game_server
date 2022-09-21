package gate

import (
	"github.com/junaozun/game_server/pkg/app"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/ws"
	"log"
)

const (
	host = "0.0.0.0:"
)

type GateApp struct {
	ServerName string
	Router     *ws.Router
}

func NewGateWay() *GateApp {
	return &GateApp{
		ServerName: "gateway",
		Router:     ws.NewRouter(),
	}
}

func (g *GateApp) Run(cfg pkgConfig.GameConfig) error {
	g.initRouter()
	wsServer := ws.NewWsServer(host+cfg.GateWay.Port, g.Router)
	gate := app.New(
		app.OnBeginHook(func() {
			log.Printf("gate app start,addr:%s ....", wsServer.Addr)
		}),
		app.OnExitHook(func() {
			log.Printf("gate app exit,addr:%s ....", wsServer.Addr)
		}),
		app.Name(g.ServerName),
		app.Runners(wsServer),
	)
	if err := gate.Run(); err != nil {
		return err
	}
	return nil
}

func (g *GateApp) initRouter() {
	g.Router.Group("*").AddRouter("*", g.routerForward)
}

func (g *GateApp) routerForward(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {
	log.Println("请求到达gateway....")
}
