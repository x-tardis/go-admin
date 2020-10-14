package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func UserCenter(v1 gin.IRouter) {
	r := v1.Group("/user")
	{
		r.GET("/profile", system.GetSysUserProfile)
		r.POST("/avatar", system.InsetSysUserAvatar)
		r.PUT("/pwd", system.SysUserUpdatePwd)
	}
}
