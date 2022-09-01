package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/web/controller/account"
	"github.com/junaozun/game_server/internal/web/repo/repo_impl_mysql"
	"github.com/junaozun/game_server/internal/web/service"
	"github.com/junaozun/game_server/pkg/httpx/middleware"
)

func RegisterRouters(component *component.Component) func(g *gin.Engine) {
	return func(g *gin.Engine) {
		g.Use(middleware.Cors())
		g.GET("/ping", func(context *gin.Context) {
			context.String(http.StatusOK, "pong")
		})

		// 用依赖注入wire替代
		repo := repo_impl_mysql.NewAccountRepo(component.Dao.DB)
		service := service.NewAccountService(repo)
		accountCtl := account.NewAccountCtl(service)

		account := g.Group("/account")
		{
			account.Any("/register", accountCtl.Register)
		}
	}
}
