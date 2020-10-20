package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/sysfile"
)

// 无需认证的路由代码
func PubSysFileInfo(v1 gin.IRouter) {
	ctl := new(sysfile.FileInfo)
	r := v1.Group("/sysfileinfo")
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
	ctl := new(sysfile.FileDir)
	r := v1.Group("/sysfiledir")
	{
		r.GET("", ctl.QueryTree)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:ids", ctl.BatchDelete)
	}
}
