package game

import (
	"log"

	"github.com/junaozun/game_server/internal/logic/component"
	"github.com/junaozun/game_server/internal/logic/model"
	"github.com/junaozun/game_server/net"
)

type Game struct {
	*MgrManager
	Component *component.Component
	Router    *net.Router
	Account   *Account
}

func NewGame(component *component.Component, router *net.Router) *Game {
	g := &Game{
		MgrManager: NewMgrManager(),
		Component:  component,
		Router:     router,
	}
	g.Init()
	return g
}

func (g *Game) Init() {
	// g.initTable()
	g.initGame()
	g.initRouter()
}

func (g *Game) initTable() {
	err := g.Component.Dao.DB.AutoMigrate(
		new(model.User),
		new(model.LoginHistory),
		new(model.LoginLast),
	)
	if err != nil {
		log.Printf("[game] initTable err:%s", err.Error())
		panic(err)
	}
}

func (g *Game) initGame() {
	// 初始化账号系统
	g.Account = NewAccount(g)

	// register
	g.Register(g.Account)

}

func (g *Game) initRouter() {
	for _, v := range g.Modules {
		command := v.RegisterRouter()
		group := g.Router.Group(command.group)
		group.AddRouter(command.action, command.execFunc)
	}
}

type ExecCommand struct {
	group    string
	action   string
	execFunc net.HandlerFunc
}
