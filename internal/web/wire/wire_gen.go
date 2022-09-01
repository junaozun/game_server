// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/junaozun/game_server/internal/web/controller/account"
	"github.com/junaozun/game_server/internal/web/repo/repo_impl_mysql"
	"github.com/junaozun/game_server/internal/web/router"
	"github.com/junaozun/game_server/internal/web/service"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func NewWebRouterMgr(db *gorm.DB) router.WebRouter {
	iAccountRepo := repo_impl_mysql.NewAccountRepo(db)
	accountService := service.NewAccountService(iAccountRepo)
	accountCtl := account.NewAccountCtl(accountService)
	webRouter := router.NewWebRouter(accountCtl)
	return webRouter
}
