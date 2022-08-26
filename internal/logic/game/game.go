package game

import (
	"log"

	"github.com/junaozun/game_server/internal/logic/component"
	"github.com/junaozun/game_server/internal/logic/model"
	"github.com/junaozun/game_server/internal/logic/wsMgr"
	"github.com/junaozun/game_server/pkg/ws"
)

type Game struct {
	*MgrManager          // mgr管理器
	*component.Component // 组件
	*ws.Router           // 路由
	*wsMgr.WsMgr         // 在线用户
	// 系统功能
	Account *Account
	// 玩法功能
}

func NewGame(component *component.Component, router *ws.Router, onLineUser *wsMgr.WsMgr) *Game {
	g := &Game{
		MgrManager: NewMgrManager(),
		Component:  component,
		Router:     router,
		WsMgr:      onLineUser,
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
		v.RegisterRouter(func(command ExecCommand) {
			g.Router.Group(command.group).AddRouter(command.name, command.execFunc)
		})
	}
}

type ExecCommand struct {
	group    string
	name     string
	execFunc ws.HandlerFunc
}
