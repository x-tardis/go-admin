package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/log"
)

func LoginLog(v1 gin.IRouter) {
	ctl := new(log.LoginLog)
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
	r := v1.Group("/operlog")
	{
		r.GET("", log.GetOperLogList)
		r.GET("/:id", log.GetOperLog)
		r.DELETE("/:ids", log.DeleteOperLog)
	}
}
