package data

import (
	"sync"
	"time"

	"github.com/junaozun/game_server/api"
)

type MapCity struct {
	mutex      sync.Mutex `gorm:"-"`
	Id         int        `gorm:"id pk autoincr"` // cityID
	RId        int        `gorm:"column:rid"`
	Name       string     `gorm:"column:name"`
	X          int        `gorm:"column:x"`
	Y          int        `gorm:"column:y"`
	IsMain     int8       `gorm:"column:isMain"`
	CurDurable int        `gorm:"column:curDurable"`
	CreatedAt  time.Time  `gorm:"column:createdAt"`
	OccupyTime time.Time  `gorm:"column:occupyTime"`
}

func (m *MapCity) TableName() string {
	return "map_city"
}

func (m *MapCity) ToClient() api.MapRoleCity {
	p := api.MapRoleCity{}
	p.X = m.X
	p.Y = m.Y
	p.CityId = m.Id
	p.UnionId = 0
	p.UnionName = ""
	p.ParentId = 0
	p.MaxDurable = 1000
	p.CurDurable = m.CurDurable
	p.Level = 1
	p.RId = m.RId
	p.Name = m.Name
	p.IsMain = m.IsMain == 1
	p.OccupyTime = m.OccupyTime.UnixNano() / 1e6
	return p
}
