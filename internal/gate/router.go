package gate

import (
	"log"

	"github.com/junaozun/game_server/pkg/ws"
)

func (g *GateApp) InitRouter() {
	g.Router.Group("*").AddRouter("*", g.routerForward)
}

func (g *GateApp) routerForward(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {
	log.Println("请求到达gateway....")
}
