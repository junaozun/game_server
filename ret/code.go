package ret

import (
	"github.com/junaozun/game_server/pkg/errno"
)

var (
	// OK
	OK                 = errno.NewError(0, "success")
	Err_UserNotFound   = errno.NewError(1001, "用户未找到")
	Err_UserExist      = errno.NewError(1002, "用户已存在")
	Err_PasswdNotRight = errno.NewError(1003, "密码不正确")

	Err_Param = errno.NewError(2000, "参数错误")
	Err_DB    = errno.NewError(2001, "数据错误")
)
