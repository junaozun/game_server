package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/global"
	"github.com/junaozun/game_server/internal/web/service"
	common_model "github.com/junaozun/game_server/model"
	"github.com/junaozun/game_server/pkg/utils"
	"github.com/junaozun/game_server/ret"
)

type AccountCtl struct {
	Service     *service.AccountService
	NatsService *service.NatsService
}

func NewAccountCtl(accountService *service.AccountService, natsService *service.NatsService) *AccountCtl {
	return &AccountCtl{
		Service:     accountService,
		NatsService: natsService,
	}
}

func (ctl *AccountCtl) Register(c *gin.Context) {
	ctx := c.Request.Context()
	req := api.RegisterReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, ret.Err_Param)
		return
	}
	dbUser := &common_model.User{
		CreateTime: global.Now(),
		Username:   req.Username,
		Passwd:     utils.ScryptPasswd(req.Password),
		Hardware:   req.Hardware,
		Status:     0,
		IsOnline:   false,
	}
	errno := ctl.Service.AddAccount(ctx, dbUser)
	c.JSON(http.StatusOK, errno)
}

func (ctl *AccountCtl) UseNatsTest(c *gin.Context) {
	ctx := c.Request.Context()
	resp, err := ctl.NatsService.UseNatsTest(ctx, "nihao")
	if err != nil {
		c.JSON(http.StatusOK, ret.Err_Param)
		return
	}
	r := &api.NatsRpcResp{
		Total: resp.Id,
		Name:  resp.Brother,
	}
	c.JSON(http.StatusOK, ret.OK.WithData(r))
}

func (ctl *AccountCtl) GetRankTest(c *gin.Context) {
	ctx := c.Request.Context()
	resp, err := ctl.NatsService.GetRankTest(ctx, "area")
	if err != nil {
		c.JSON(http.StatusOK, ret.Err_Param)
		return
	}
	c.JSON(http.StatusOK, ret.OK.WithData(resp))
}
