package data

import (
	"time"

	"github.com/junaozun/game_server/api"
)

type MapRoleBuild struct {
	Id         int       `gorm:"id pk autoincr"`
	RId        int       `gorm:"column:rid"`
	Type       int8      `gorm:"column:type"`
	Level      int8      `gorm:"column:level"`
	OPLevel    int8      `gorm:"column:opLevel"` // 操作level
	X          int       `gorm:"column:x"`
	Y          int       `gorm:"column:y"`
	Name       string    `gorm:"column:name"`
	Wood       int       `gorm:"-"`
	Iron       int       `gorm:"-"`
	Stone      int       `gorm:"-"`
	Grain      int       `gorm:"-"`
	Defender   int       `gorm:"-"`
	CurDurable int       `gorm:"column:curDurable"`
	MaxDurable int       `gorm:"column:maxDurable"`
	OccupyTime time.Time `gorm:"column:occupyTime"`
	EndTime    time.Time `gorm:"column:endTime"` // 建造或升级完的时间
	GiveUpTime int64     `gorm:"column:giveUpTime"`
}

func (m *MapRoleBuild) TableName() string {
	return "map_role_build"
}

func (m *MapRoleBuild) ToClient() api.MapRoleBuild {
	p := api.MapRoleBuild{}
	p.RNick = "111"
	p.UnionId = 0
	p.UnionName = ""
	p.ParentId = 0
	p.X = m.X
	p.Y = m.Y
	p.Type = m.Type
	p.RId = m.RId
	p.Name = m.Name

	p.OccupyTime = m.OccupyTime.UnixNano() / 1e6
	p.GiveUpTime = m.GiveUpTime * 1000
	p.EndTime = m.EndTime.UnixNano() / 1e6

	p.CurDurable = m.CurDurable
	p.MaxDurable = m.MaxDurable
	p.Defender = m.Defender
	p.Level = m.Level
	p.OPLevel = m.OPLevel
	return p
}
