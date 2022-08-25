package logic

import (
	"context"
	"flag"
	"github.com/junaozun/game_server/internal/logic/component"
	"github.com/junaozun/game_server/internal/logic/game"
	"github.com/junaozun/game_server/net"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/db"
	"log"
)

var (
	ServerId   string
	ServerPort string
)

const (
	host = "0.0.0.0:"
)

type LogicService struct {
	wsServer   *net.Server
	component  *component.Component
	router     *net.Router
	onLineUser *net.WsMgr
}

func NewLogicService() *LogicService {
	logicService := &LogicService{}
	return logicService
}

func (l *LogicService) ParseFlag(set *flag.FlagSet) {
	set.StringVar(&ServerId, "server_id", "", "logic server id")
}

func (l *LogicService) Init(cfg pkgConfig.GameConfig) error {
	// 初始化数据访问层
	dao, err := db.NewDao(cfg.Logic.Mysql)
	if err != nil {
		return err
	}
	ServerPort = cfg.Logic.Port
	l.component = component.NewComponent(dao)
	l.router = net.NewRouter()
	l.onLineUser = net.NewWsMgr()
	// 初始化游戏玩法
	game.NewGame(l.component, l.router, l.onLineUser)
	log.Println("[LogicService] init successful.....")
	return nil
}

func (l *LogicService) Start(ctx context.Context) error {
	log.Println("[LogicService] start........")
	l.wsServer = net.NewServer(host+ServerPort, l.router)
	return l.wsServer.Start(ctx)
}

func (l *LogicService) Stop(ctx context.Context) error {
	log.Println("[LogicService] stop ........")
	return l.wsServer.Stop(ctx)
}
