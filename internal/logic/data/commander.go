package data

import (
	"time"

	"github.com/junaozun/game_server/api"
	common_model "github.com/junaozun/game_server/model"
)

const (
	CommanderNormal      = 0 // 正常
	CommanderComposeStar = 1 // 星级合成
	CommanderConvert     = 2 // 转换
)
const SkillLimit = 3

// Commander 武将
type Commander struct {
	Id            int                            `gorm:"id pk autoincr"`
	RId           int                            `gorm:"column:rid"`
	CfgId         int                            `gorm:"column:cfgId"`
	PhysicalPower int                            `gorm:"column:physicalPower"`
	Level         int8                           `gorm:"column:level"`
	Exp           int                            `gorm:"column:exp"`
	Order         int8                           `gorm:"column:order"`
	CityId        int                            `gorm:"column:cityId"`
	CreatedAt     time.Time                      `gorm:"column:createdAt"`
	CurArms       int                            `gorm:"column:arms"`
	HasPrPoint    int                            `gorm:"column:hasPrPoint"`
	UsePrPoint    int                            `gorm:"column:usePrPoint"`
	AttackDis     int                            `gorm:"column:attackDistance"`
	ForceAdded    int                            `gorm:"column:forceAdded"`
	StrategyAdded int                            `gorm:"column:strategyAdded"`
	DefenseAdded  int                            `gorm:"column:defenseAdded"`
	SpeedAdded    int                            `gorm:"column:speedAdded"`
	DestroyAdded  int                            `gorm:"column:destroyAdded"`
	StarLv        int8                           `gorm:"column:starLv"`
	Star          int8                           `gorm:"column:star"`
	ParentId      int                            `gorm:"column:parentId"`
	Skills        string                         `gorm:"column:skills"`
	SkillsArray   []*common_model.CommanderSkill `gorm:"-"`
	State         int8                           `gorm:"column:state"`
}

func (g *Commander) TableName() string {
	return "commander"
}

func (g *Commander) ToClient() api.Commander {
	p := api.Commander{}
	p.CityId = g.CityId
	p.Order = g.Order
	p.PhysicalPower = g.PhysicalPower
	p.Id = g.Id
	p.CfgId = g.CfgId
	p.Level = g.Level
	p.Exp = g.Exp
	p.CurArms = g.CurArms
	p.HasPrPoint = g.HasPrPoint
	p.UsePrPoint = g.UsePrPoint
	p.AttackDis = g.AttackDis
	p.ForceAdded = g.ForceAdded
	p.StrategyAdded = g.StrategyAdded
	p.DefenseAdded = g.DefenseAdded
	p.SpeedAdded = g.SpeedAdded
	p.DestroyAdded = g.DestroyAdded
	p.StarLv = g.StarLv
	p.Star = g.Star
	p.State = g.State
	p.ParentId = g.ParentId
	p.Skills = g.SkillsArray
	return p
}
