package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func User(v1 gin.IRouter) {
	// v1.GET("/sysUser", system.GetSysUserInit)
	r := v1.Group("/users")
	{
		r.GET("", system.GetSysUserList)
		r.GET("/:id", system.GetSysUser)
		r.GET("/", system.GetSysUserInit)
		r.POST("", system.InsertSysUser)
		r.PUT("", system.UpdateSysUser)
		r.DELETE("/:ids", system.DeleteSysUser)
	}
}
