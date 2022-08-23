package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/junaozun/game_server/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			// 将标准输出作为Writer
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				// 设定慢查询时间阈值为10ms
				SlowThreshold: 10 * time.Microsecond,
				// 设置日志级别，只有Warn和Info级别会输出慢查询日志
				LogLevel: logger.Warn,
			},
		),
	})
	if nil != err {
		return nil, err
	}
	dao := &Dao{
		DB: db,
	}
	return dao, nil
}
