package model

import (
	"time"
)

// User 用户表
type User struct {
	Id         int       `gorm:"primaryKey;column:id;type:int(10) unsigned;not null"`
	UId        int       `gorm:"column:uid;type:int(20)"`
	Username   string    `gorm:"column:username;type:varchar(20)" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	Passcode   string    `gorm:"column:passcode;type:varchar(11)"`
	Passwd     string    `gorm:"column:passwd;type:varchar(20)"`
	Hardware   string    `gorm:"column:hardware;type:varchar(11)"`
	Status     int       `gorm:"column:status;type:int(11)"`
	CreateTime time.Time `gorm:"column:createTime;type:datetime"`
	IsOnline   bool      `gorm:"-"`
}

// func (*User) TableName() string {
// 	return "user"
// }
