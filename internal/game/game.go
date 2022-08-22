package game

import (
	"github.com/junaozun/game_server/net"
)

type Game struct {
	ServerMgr *net.Server
	Account   *Account
}

func NewGame(s *net.Server) *Game {
	g := &Game{
		ServerMgr: s,
	}
	g.Init()
	return g
}

func (g *Game) Init() {
	// 初始化table
	g.InitTable()
	// 初始化路由
	g.InitGame()
}

func (g *Game) InitTable() {
	g.ServerMgr.DBEngine.AutoMigrate()
}

func (g *Game) InitGame() {
	// 账号
	g.Account = NewAccount(g)
	g.Account.RegisterRouter()

	// xx
}
