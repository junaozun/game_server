package component

import (
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/dao"
	"github.com/junaozun/game_server/pkg/natsx"
)

// Component 组件
type Component struct {
	Dao         *dao.Dao      // 数据访问层组件（mysql,redis,etcd）
	LogicClient *natsx.Client // nats 消息中间件组件
	ChessClient *natsx.Client // nats 消息中间件组件
	GvgClient   *natsx.Client // nats 消息中间件组件
	// kafka // 消息中间件组件
}

func NewComponent(dao *dao.Dao, cfg pkgConfig.GameConfig) *Component {
	connEnc, err := natsx.NewNatsJSONEnc(cfg.Common.NATS.Server)
	if err != nil {
		panic(err)
	}
	logicClient, err := natsx.NewClient(connEnc, "logic")
	if err != nil {
		panic(err)
	}
	chessClient, err := natsx.NewClient(connEnc, "chess")
	if err != nil {
		panic(err)
	}
	gvgClient, err := natsx.NewClient(connEnc, "gvg")
	if err != nil {
		panic(err)
	}
	return &Component{
		Dao:         dao,
		LogicClient: logicClient,
		ChessClient: chessClient,
		GvgClient:   gvgClient,
	}
}
