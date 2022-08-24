package logic

import (
	"flag"

	"github.com/junaozun/game_server/internal/logic/component"
	"github.com/junaozun/game_server/internal/logic/game"
	"github.com/junaozun/game_server/net"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/db"
)

const (
	host = "0.0.0.0:"
)

type LogicService struct {
	serverPort string
	serverId   string
	component  *component.Component
	router     *net.Router
	closeChan  chan struct{}
	fc         chan func()
}

func NewLogicService() *LogicService {
	logicService := &LogicService{
		closeChan: make(chan struct{}),
		fc:        make(chan func(), 32),
	}
	return logicService
}

func (l *LogicService) ParseFlag(set *flag.FlagSet) {
	set.StringVar(&l.serverId, "server_id", "", "logic server id")
}

func (l *LogicService) Init(cfg pkgConfig.GameConfig) error {
	// 初始化数据访问层
	dao, err := db.NewDao(cfg.Logic.Mysql)
	if err != nil {
		return err
	}
	l.serverPort = cfg.Logic.Port
	l.component = component.NewComponent(dao)
	l.router = net.NewRouter()
	// 初始化游戏
	game.NewGame(l.component, l.router)
	return nil
}

func (l *LogicService) Run() {
	server := net.NewServer(host+l.serverPort, l.router)
	server.Start()
}
