package data

import (
	"time"

	"github.com/junaozun/game_server/api"
)

// 角色信息表

type Role struct {
	Id         int       `gorm:"rid pk autoincr"` // 角色id
	UId        int       `gorm:"column:uid"`
	NickName   string    `gorm:"column:nickName;type:varchar(50)"`
	Balance    int       `gorm:"column:balance"`
	HeadId     int16     `gorm:"column:headId"`
	Sex        int8      `gorm:"column:sex"`
	Profile    string    `gorm:"column:profile;type:varchar(30)"`
	LoginTime  time.Time `gorm:"column:loginTime;DATETIME(10)"`
	LogoutTime time.Time `gorm:"column:logoutTime;DATETIME(10)"`
	CreatedAt  time.Time `gorm:"column:createdAt;DATETIME(10)"`
}

func (r *Role) TableName() string {
	return "role"
}

func (r *Role) ToClient() api.RoleBase {
	m := api.RoleBase{}
	m.UId = r.UId
	m.RId = r.Id
	m.Sex = r.Sex
	m.NickName = r.NickName
	m.HeadId = r.HeadId
	m.Balance = r.Balance
	m.Profile = r.Profile
	return m
}
