package logic

import (
	"context"
	"flag"
	"log"

	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/logic/game"
	"github.com/junaozun/game_server/internal/logic/wsMgr"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/dao"
	"github.com/junaozun/game_server/pkg/ws"
)

var (
	ServerId   string
	ServerPort string
)

const (
	host = "0.0.0.0:"
)

type LogicService struct {
	wsServer   *ws.WsServer
	component  *component.Component
	onLineUser *wsMgr.WsMgr
}

func NewLogicService() *LogicService {
	logicService := &LogicService{
		onLineUser: wsMgr.NewWsMgr(),
	}
	return logicService
}

func (l *LogicService) ParseFlag(set *flag.FlagSet) {
	set.StringVar(&ServerId, "server_id", "", "logic server id")
}

func (l *LogicService) Init(cfg pkgConfig.GameConfig) error {
	// 初始化数据访问层
	dao, err := dao.NewDao([]interface{}{cfg.Logic.Mysql, cfg.Common.Etcd, cfg.Common.Cache})
	if err != nil {
		return err
	}
	ServerPort = cfg.Logic.Port
	l.wsServer = ws.NewWsServer(host+ServerPort, ws.NewRouter())
	l.component = component.NewComponent(dao)
	// 初始化游戏玩法
	game.NewGame(l.component, l.wsServer.ServerRouter, l.onLineUser)
	log.Println("[LogicService] init successful.....")
	return nil
}

func (l *LogicService) Start(ctx context.Context) error {
	log.Printf("[LogicService] addr:%s,start........", host+ServerPort)
	return l.wsServer.Start(ctx)
}

func (l *LogicService) Stop(ctx context.Context) error {
	log.Printf("[LogicService] addr:%s stop ........", host+ServerPort)
	return l.wsServer.Stop(ctx)
}
