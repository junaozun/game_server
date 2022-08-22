package game

import (
	"log"

	"github.com/junaozun/game_server/internal/model"
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
	err := g.ServerMgr.DBEngine.AutoMigrate(
		new(model.User),
		new(model.LoginHistory),
		new(model.LoginLast),
	)
	if err != nil {
		log.Printf("[game] initTable err:%s", err.Error())
		panic(err)
	}
}

func (g *Game) InitGame() {
	// 账号
	g.Account = NewAccount(g)
	g.Account.RegisterRouter()

	// xx
}
