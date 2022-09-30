package gate

import (
	"log"
	"strings"

	"github.com/junaozun/game_server/pkg/ws"
)

func (g *GateApp) InitRouter() {
	g.Router.Group("*").AddRouter("*", g.routerForward)
}

func (g *GateApp) routerForward(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {
	log.Println("请求到达gateway....")
	routerName := req.Body.Router
	var proxyAddr string
	if isAccount(routerName) {
		proxyAddr = g.Handler.GetLoginProxy()
	} else {
		proxyAddr = g.Handler.GetLogicProxy()
	}
	proxyClient := NewProxyClient(proxyAddr)
	proxyClient.Connect()
}

func isAccount(routerName string) bool {
	return strings.HasPrefix(routerName, "account.")
}
