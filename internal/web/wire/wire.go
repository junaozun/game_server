//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/web/controller/account"
	"github.com/junaozun/game_server/internal/web/repo/repo_impl_mysql"
	"github.com/junaozun/game_server/internal/web/repo/repo_impl_nats"
	"github.com/junaozun/game_server/internal/web/router"
	"github.com/junaozun/game_server/internal/web/service"
)

func NewWebRouterMgr(commponent *component.Component) router.WebRouter {
	wire.Build(
		// repo
		repo_impl_mysql.NewAccountRepo,
		repo_impl_nats.NewNatsRepo,

		// service
		service.NewAccountService,
		service.NewNatsService,

		// controller
		account.NewAccountCtl,

		// router
		router.NewWebRouter,
	)
	return nil
}
