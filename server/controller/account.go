package controller

import (
	"github.com/junaozun/game_server/net"
	"github.com/junaozun/game_server/server/proto"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(router *net.Router) {
	g := router.Group("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	rsp.Body.Code = 0
	loginResp := &proto.LoginRsp{
		Username: "admin",
		Session:  "nnnnnn",
		UId:      1,
	}
	rsp.Body.Msg = loginResp
}
