package game

import (
	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/net"
)

type Account struct {
	g *Game
}

func NewAccount(g *Game) *Account {
	return &Account{
		g: g,
	}
}

func (a *Account) RegisterRouter() {
	g := a.g.ServerMgr.ServerRouter.Group("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	// 1.拿到前端的用户名和密码和硬件Id

	// 2.根据用户名查询user，得到用户数据
	// 3.进行密码比对，如果密码正确 登录成功
	// 4. 保存用户登录记录
	// 5. 保存用户的最后一次登录
	// 6.客户端 需要一个session (JWT生成)
	rsp.Body.Code = 0
	loginResp := &api.LoginRsp{
		Username: "admin",
		Session:  "nnnnnn",
		UId:      1,
	}
	rsp.Body.Msg = loginResp
}
