package db

import (
	"fmt"

	"github.com/junaozun/game_server/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Dao 数据对象访问
type Dao struct {
	DB *gorm.DB // db
}

// NewDao 构造
func NewDao(daoConfig config.DBConfig) (*Dao, error) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		daoConfig.User,
		daoConfig.Password,
		daoConfig.Host,
		daoConfig.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if nil != err {
		return nil, err
	}
	dao := &Dao{
		DB: db,
	}
	return dao, nil
}
