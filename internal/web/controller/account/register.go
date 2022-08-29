package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/global"
	common_model "github.com/junaozun/game_server/model"
	"github.com/junaozun/game_server/pkg/dao"
	"github.com/junaozun/game_server/pkg/utils"
)

type RegisterAccount struct {
	Dao *dao.Dao
}

func NewAccount(dao *dao.Dao) *RegisterAccount {
	return &RegisterAccount{
		Dao: dao,
	}
}

func (ctl *RegisterAccount) Register(c *gin.Context) {
	req := api.RegisterReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, &RET{
			Code: -1,
			Msg:  "shouldBind err",
			Data: nil,
		})
		return
	}
	user := &common_model.User{}
	err := ctl.Dao.DB.Where(&common_model.User{Username: req.Username}).Limit(1).Find(user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, &RET{
			Code: -2,
			Msg:  "db err",
			Data: nil,
		})
		return
	}
	if user.UId != 0 {
		c.JSON(http.StatusInternalServerError, &RET{
			Code: -3,
			Msg:  "用户名已存在",
			Data: nil,
		})
		return
	}
	dbUser := &common_model.User{
		CreateTime: global.Now(),
		UId:        1,
		Username:   req.Username,
		Passwd:     utils.ScryptPasswd(req.Password),
		Hardware:   req.Hardware,
		Status:     0,
		IsOnline:   false,
	}
	err = ctl.Dao.DB.Create(dbUser).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, &RET{
			Code: -4,
			Msg:  "insert user err",
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, &RET{
		Code: 0,
		Msg:  "success",
		Data: nil,
	})
}

type RET struct {
	Code int         `json:"code"`    // 业务code 0表示成功,非0表示失败
	Msg  string      `json:"message"` // 错误描述 成功是OK
	Data interface{} `json:"data"`    // 成功时返回的数据
}
