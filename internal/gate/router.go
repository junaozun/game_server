package gate

import (
	"strings"
	"time"

	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/game_server/ret"
	"github.com/junaozun/gogopkg/logrusx"
	"github.com/mitchellh/mapstructure"
)

func (g *GateApp) InitRouter() {
	g.Router.Group("*").AddRouter("*", g.routerForward)
}

func (g *GateApp) routerForward(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {
	logrusx.Log.WithFields(logrusx.Fields{}).Info("[GateApp] routerForward 客户端请求到达gateway....")
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Router = req.Body.Router
	rsp.Body.Code = ret.OK.Code

	routerName := req.Body.Router

	if routerName == "" {
		rsp.Body.Code = ret.Err_ProxyNotFound.Code
		return
	}
	// 心跳消息包直接由gateway回
	if routerName == common.HearbeatMsg {
		h := &ws.Hearbeat{}
		mapstructure.Decode(req.Body.Msg, h)
		h.ServerTime = time.Now().UnixNano() / 1e6
		rsp.Body.Msg = h
		return
	}
	var proxyAddr string
	// if isAccount(routerName) {
	// 	proxyAddr = g.Handler.GetLoginProxy()
	// } else {
	proxyAddr = g.Handler.GetLogicProxy()
	// }

	// 客户端id
	c, ok := req.Conn.GetProperty("cid")
	if !ok {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[GateApp] routerForward cid not found")
		rsp.Body.Code = ret.Err_Param.Code
		return
	}
	cid := c.(int64)

	proxyClients := g.Handler.GetProxyMap(proxyAddr)
	// 获取该客户端与logic或login的ws连接
	proxyClient, ok := proxyClients[cid]
	// 没有找到，则新建立链接
	if !ok {
		proxyClient = NewProxyClient(proxyAddr)
		err := proxyClient.ConnectServer()
		if err != nil {
			logrusx.Log.WithFields(logrusx.Fields{}).Error("[GateApp] routerForward proxyClient.ConnectServer err")
			rsp.Body.Code = ret.Err_ProxyConnect.Code
			return
		}
		g.Handler.SetProxyMapValue(proxyAddr, cid, proxyClient)
		// 给链接中设置属性
		proxyClient.SetProperty("cid", cid)
		proxyClient.SetProperty("proxyAddr", proxyAddr)
		proxyClient.SetProperty("clientConn", req.Conn)
		proxyClient.OnPushClient(g.Handler.OnPushClient)
	}
	// 给logic server 发送数据,并等待数据返回
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
