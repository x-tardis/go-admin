package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func SysUser(v1 gin.IRouter) {
	v1.GET("/sysUserList", system.GetSysUserList)
	r := v1.Group("/sysUser")
	{
		r.GET("/:userId", system.GetSysUser)
		r.GET("/", system.GetSysUserInit)
		r.POST("", system.InsertSysUser)
		r.PUT("", system.UpdateSysUser)
		r.DELETE("/:userId", system.DeleteSysUser)
	}
}
