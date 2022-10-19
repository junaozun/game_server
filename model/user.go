package model

import (
	"time"
)

// User 用户表
type User struct {
	ID         uint64 `gorm:"primarykey"` // uid
	Username   string `gorm:"column:username; unique type:varchar(20)" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	Passwd     string `gorm:"column:passwd;type:varchar(100)"`
	Hardware   string `gorm:"column:hardware;type:varchar(100)"`
	CreateTime time.Time
	Status     int  `gorm:"column:status;type:int(11)"`
	IsOnline   bool `gorm:"-"`
}

func (*User) TableName() string {
	return "user"
}
