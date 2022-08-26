package component

import (
	"github.com/junaozun/game_server/pkg/dao"
)

// Component 组件
type Component struct {
	Dao *dao.Dao // 数据访问层组件（mysql,redis,etcd）
	// nats // 消息中间件组件
	// kafka // 消息中间件组件
}

func NewComponent(dao *dao.Dao) *Component {
	return &Component{
		Dao: dao,
	}
}
