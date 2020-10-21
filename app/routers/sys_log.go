package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func LoginLog(v1 gin.IRouter) {
	ctl := new(system.LoginLog)
	r := v1.Group("/loginlog")
	{
		r.GET("", ctl.QueryPage)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:ids", ctl.BatchDelete)
	}
}

func OperLog(v1 gin.IRouter) {
	ctl := new(system.OperLog)

	r := v1.Group("/operlog")
	{
		r.GET("", ctl.QueryPage)
		r.GET("/:id", ctl.Get)
		r.DELETE("/:ids", ctl.BatchDelete)
	}
}
