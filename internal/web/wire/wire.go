//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/junaozun/game_server/internal/web/controller/account"
	"github.com/junaozun/game_server/internal/web/repo/repo_impl_mysql"
	"github.com/junaozun/game_server/internal/web/router"
	"github.com/junaozun/game_server/internal/web/service"
	"gorm.io/gorm"
)

func NewWebRouterMgr(db *gorm.DB) router.WebRouter {
	wire.Build(
		// repo
		repo_impl_mysql.NewAccountRepo,
		// service
		service.NewAccountService,
		// controller
		account.NewAccountCtl,

		// router
		router.NewWebRouter,
	)
	return nil
}
