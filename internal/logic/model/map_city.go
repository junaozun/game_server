package model

import (
	"sync"
	"time"
)

type MapCity struct {
	mutex      sync.Mutex `gorm:"-"`
	CityId     int        `gorm:"column:cityId;pk autoincr"`
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

// func (m *MapRoleCity) ToModel() interface{} {
// 	p := model.MapRoleCity{}
// 	p.X = m.X
// 	p.Y = m.Y
// 	p.CityId = m.CityId
// 	p.UnionId = GetUnion(m.RId)
// 	p.UnionName = ""
// 	p.ParentId = 0
// 	p.MaxDurable = 1000
// 	p.CurDurable = m.CurDurable
// 	p.Level = 1
// 	p.RId = m.RId
// 	p.Name = m.Name
// 	p.IsMain = m.IsMain == 1
// 	p.OccupyTime = m.OccupyTime.UnixNano() / 1e6
// 	return p
// }
