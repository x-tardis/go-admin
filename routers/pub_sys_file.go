package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/system"
)

// 无需认证的路由代码
func PubSysFileInfo(v1 gin.IRouter) {
	ctl := new(system.FileInfo)
	r := v1.Group("/fileinfo")
	{
		r.GET("", ctl.QueryPage)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:ids", ctl.BatchDelete)
	}
}

// 无需认证的路由代码
func PubSysFileDir(v1 gin.IRouter) {
	ctl := new(system.FileDir)
	r := v1.Group("/filedir")
	{
		r.GET("", ctl.QueryTree)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:ids", ctl.BatchDelete)
	}
}
