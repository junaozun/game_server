package model

import (
	"github.com/junaozun/game_server/internal/logic"
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	UId      int    `gorm:"column:uid;type:int(20)"`
	Username string `gorm:"column:username;type:varchar(20)" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	Passcode string `gorm:"column:passcode;type:varchar(11)"`
	Passwd   string `gorm:"column:passwd;type:varchar(20)"`
	Hardware string `gorm:"column:hardware;type:varchar(11)"`
	Status   int    `gorm:"column:status;type:int(11)"`
	IsOnline bool   `gorm:"-"`
}

func (*User) TableName() string {
	return "user_" + logic.ServerId
}
