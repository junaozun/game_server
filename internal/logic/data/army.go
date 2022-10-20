package data

import (
	"time"

	"github.com/junaozun/game_server/api"
)

const (
	ArmyCmdIdle        = 0 // 空闲
	ArmyCmdAttack      = 1 // 攻击
	ArmyCmdDefend      = 2 // 驻守
	ArmyCmdReclamation = 3 // 屯垦
	ArmyCmdBack        = 4 // 撤退
	ArmyCmdConscript   = 5 // 征兵
	ArmyCmdTransfer    = 6 // 调动
)

const (
	ArmyStop    = 0
	ArmyRunning = 1
)

// Army 军队
type Army struct {
	Id                 int          `gorm:"id pk autoincr"`
	RId                int          `gorm:"column:rid"`
	CityId             int          `gorm:"column:cityId"`
	Order              int8         `gorm:"column:order"`
	Commanders         string       `gorm:"column:commanders"`
	Soldiers           string       `gorm:"column:soldiers"`
	ConscriptTimes     string       `gorm:"column:conscriptTimes"` // 征兵结束时间，json数组
	ConscriptCnts      string       `gorm:"column:conscriptCnts"`  // 征兵数量，json数组
	Cmd                int8         `gorm:"column:cmd"`
	FromX              int          `gorm:"column:fromX"`
	FromY              int          `gorm:"column:fromY"`
	ToX                int          `gorm:"column:toX"`
	ToY                int          `gorm:"column:toY"`
	Start              time.Time    `gorm:"column:start"`
	End                time.Time    `gorm:"column:end"`
	State              int8         `gorm:"-"` // 状态:0:running,1:stop
	GeneralArray       []int        `gorm:"-"`
	SoldierArray       []int        `gorm:"-"`
	ConscriptTimeArray []int64      `gorm:"-"`
	ConscriptCntArray  []int        `gorm:"-"`
	CommandersConvert  []*Commander `gorm:"-"` // 转换后的结构
	CellX              int          `gorm:"-"`
	CellY              int          `gorm:"-"`
}

func (a *Army) TableName() string {
	return "army"
}

func (a *Army) ToClient() api.Army {
	p := api.Army{}
	p.CityId = a.CityId
	p.Id = a.Id
	p.UnionId = 0
	p.Order = a.Order
	p.Generals = a.GeneralArray
	p.Soldiers = a.SoldierArray
	p.ConTimes = a.ConscriptTimeArray
	p.ConCnts = a.ConscriptCntArray
	p.Cmd = a.Cmd
	p.State = a.State
	p.FromX = a.FromX
	p.FromY = a.FromY
	p.ToX = a.ToX
	p.ToY = a.ToY
	p.Start = a.Start.Unix()
	p.End = a.End.Unix()
	return p
}
