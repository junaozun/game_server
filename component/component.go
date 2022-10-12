package component

import (
	"github.com/junaozun/game_server/common"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/dao"
	"github.com/junaozun/gogopkg/natsx"
)

// Component 组件
type Component struct {
	Dao        *dao.Dao                 // 数据访问层组件（mysql,redis,etcd）
	NatsClient map[string]*natsx.Client // key:serverName   vale:natsClient
	// kafka // 消息中间件组件
}

func NewComponent(dao *dao.Dao, cfg config.GameConfig) *Component {
	connEnc, err := natsx.NewNatsPBEnc(cfg.Common.NATS.Server)
	if err != nil {
		panic(err)
	}
	component := &Component{
		Dao:        dao,
		NatsClient: make(map[string]*natsx.Client),
	}
	for _, serverName := range common.ServerNames {
		client, err := natsx.NewClient(connEnc, serverName)
		if err != nil {
			panic(err)
		}
		component.NatsClient[serverName] = client
	}
	return component
}
