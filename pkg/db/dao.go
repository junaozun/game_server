package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/junaozun/game_server/pkg/config"
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
	db, err := gorm.Open("mysql", dsn)
	if nil != err {
		return nil, err
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(daoConfig.Debug)

	dao := &Dao{
		DB: db,
	}
	return dao, nil
}

func (d *Dao) Close() {
	d.DB.Close()
}
