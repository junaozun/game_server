package gate

import (
	"log"
	"strings"

	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/game_server/ret"
	"github.com/junaozun/gogopkg/logrusx"
)

func (g *GateApp) InitRouter() {
	g.Router.Group("*").AddRouter("*", g.routerForward)
}

func (g *GateApp) routerForward(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {
	log.Println("客户端请求到达gateway....")
	routerName := req.Body.Router
	var proxyAddr string
	if isAccount(routerName) {
		proxyAddr = g.Handler.GetLoginProxy()
	} else {
		proxyAddr = g.Handler.GetLogicProxy()
	}
	if proxyAddr == "" {
		rsp.Body.Code = ret.Err_ProxyNotFound.Code
		return
	}

	// 客户端id
	c, ok := req.Conn.GetProperty("cid")
	if !ok {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[GateApp] routerForward cid not found")
		rsp.Body.Code = ret.Err_Param.Code
		return
	}
	cid := c.(int)

	mapClient := g.Handler.GetProxyMap(proxyAddr)
	proxyClient, ok := mapClient[cid]
	if !ok {
		proxyClient = NewProxyClient(proxyAddr)
		err := proxyClient.ConnectServer()
		if err != nil {
			g.Handler.DeleteCid(proxyAddr, cid)
			logrusx.Log.WithFields(logrusx.Fields{}).Error("[GateApp] routerForward proxyClient.ConnectServer err")
			rsp.Body.Code = ret.Err_ProxyConnect.Code
			return
		}
		g.Handler.SetProxyMapValue(proxyAddr, cid, proxyClient)
		// 给链接中设置属性
		proxyClient.SetProperty("cid", cid)
		proxyClient.SetProperty("proxyAddr", proxyAddr)
		proxyClient.SetProperty("clientConn", req.Conn)
		proxyClient.OnPush(g.Handler.OnPush)
	}
	// 给loginc or login server 发送数据
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Router = req.Body.Router
	r, err := proxyClient.Send(req.Body.Router, req.Body.Msg)
	if err != nil {
		rsp.Body.Code = ret.Err_ProxyConnect.Code
		return
	}
	rsp.Body.Code = r.Code
	rsp.Body.Msg = r.Msg
}

func isAccount(routerName string) bool {
	return strings.HasPrefix(routerName, "account.")
}
