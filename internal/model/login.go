package model

import (
	"time"
)

const (
	UserStatus_Login  = 0
	UserStatus_LogOut = 1
)

// LoginHistory 登录历史记录表
type LoginHistory struct {
	Id        int       `gorm:"primaryKey;column:id;type:int(10) unsigned;not null"`
	UId       int       `gorm:"column:uid;type:int(11)"`
	UserName  string    `gorm:"column:serverId;type:varchar(11)"`
	LoginTime time.Time `gorm:"column:loginTime;type:datetime"`
	Ip        string    `gorm:"column:ip;type:varchar(11)"`
	State     int8      `gorm:"column:state;type:int(11)"`
	Hardware  string    `gorm:"column:hardware;type:varchar(11)"`
}

type LoginLast struct {
	Id         int       `gorm:"primaryKey;column:id;type:int(10) unsigned;not null"`
	UId        int       `gorm:"column:uid;type:int(11)"`
	LoginTime  time.Time `gorm:"column:loginTime;type:datetime"`
	LogoutTime time.Time `gorm:"column:logoutTime;type:datetime"`
	Ip         string    `gorm:"column:ip;type:varchar(11)"`
	Session    string    `gorm:"column:session;type:varchar(255)"`
	IsLogout   int8      `gorm:"column:isLogout;type:int(11)"`
	Hardware   string    `gorm:"column:hardware;type:varchar(11)"`
}

// func (*LoginHistory) TableName() string {
// 	return "login_history"
// }
//
// func (*LoginLast) TableName() string {
// 	return "login_last"
// }
