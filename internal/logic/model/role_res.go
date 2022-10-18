package model

import (
	"github.com/junaozun/game_server/api"
)

// 角色资源表

type RoleRes struct {
	Id     int `gorm:"id pk autoincr"`
	RId    int `gorm:"column:rid"`
	Wood   int `gorm:"column:wood"`
	Iron   int `gorm:"column:iron"`
	Stone  int `gorm:"column:stone"`
	Grain  int `gorm:"column:grain"`
	Gold   int `gorm:"column:gold"`
	Decree int `gorm:"column:decree"` // 令牌
}

func (r *RoleRes) TableName() string {
	return "role_res"
}

func (r *RoleRes) ToClient() api.RoleRes {
	m := api.RoleRes{
		Wood:          r.Wood,
		Iron:          r.Iron,
		Stone:         r.Stone,
		Grain:         r.Grain,
		Gold:          r.Gold,
		Decree:        r.Decree,
		WoodYield:     100,
		IronYield:     100,
		StoneYield:    100,
		GrainYield:    100,
		GoldYield:     100,
		DepotCapacity: 10000,
	}
	return m
}
