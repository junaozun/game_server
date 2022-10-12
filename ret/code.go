package ret

import (
	"github.com/junaozun/gogopkg/errno"
)

var (
	// OK
	OK                 = errno.NewError(0, "success")
	Err_UserNotFound   = errno.NewError(1001, "用户未找到")
	Err_UserExist      = errno.NewError(1002, "用户已存在")
	Err_PasswdNotRight = errno.NewError(1003, "密码不正确")
	Err_ProxyNotFound  = errno.NewError(1100, "代理服务未发现")
	Err_ProxyConnect   = errno.NewError(1101, "代理链接错误")

	Err_Param = errno.NewError(2000, "参数错误")
	Err_DB    = errno.NewError(2001, "数据错误")
)
