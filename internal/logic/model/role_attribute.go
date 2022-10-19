package model

import (
	"time"

	"github.com/junaozun/game_server/model"
)

type RoleAttribute struct {
	Id              int            `gorm:"id pk autoincr"`
	RId             int            `gorm:"column:rid"`
	UnionId         int            `gorm:"-"`                      // 联盟id
	ParentId        int            `gorm:"column:parentId"`        // 上级id（被沦陷）
	CollectTimes    int8           `gorm:"column:collectTimes"`    // 征收次数
	LastCollectTime time.Time      `gorm:"column:lastCollectTime"` // 最后征收的时间
	PosTags         string         `gorm:"column:posTags"`         // 位置标记
	PosTagArray     []model.PosTag `gorm:"-"`
}

func (r *RoleAttribute) TableName() string {
	return "role_attribute"
}
