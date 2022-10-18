package game

import (
	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/logic/game_config"
	"github.com/junaozun/game_server/internal/logic/model"
	"github.com/junaozun/game_server/pkg/ws"
)

type Game struct {
	*MgrManager          // mgr管理器
	*component.Component // 组件
	*ws.Router           // 路由

	// 系统功能
	Role *Role
	// 玩法功能
}

func NewGame(component *component.Component, router *ws.Router) *Game {
	g := &Game{
		MgrManager: NewMgrManager(),
		Component:  component,
		Router:     router,
	}
	g.Init()
	return g
}

func (g *Game) Init() {
	g.initTable()
	g.initGame()
	g.initRouter()
}

func (g *Game) initTable() {
	err := g.Component.Dao.DB.AutoMigrate(
		new(model.Role),
		new(model.RoleRes),
	)
	if err != nil {
		panic(err)
	}
}

func (g *Game) initGame() {

	// 初始化角色资源
	game_config.Base.Load()

	// 初始化角色
	g.Role = NewRole(g)

	// register
	g.Register(g.Role)

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
