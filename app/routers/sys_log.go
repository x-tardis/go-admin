package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/log"
)

func LoginLog(v1 gin.IRouter) {
	r := v1.Group("/loginlog")
	{
		r.GET("", log.GetLoginLogList)
		r.GET("/:id", log.GetLoginLog)
		r.POST("", log.InsertLoginLog)
		r.PUT("", log.UpdateLoginLog)
		r.DELETE("/:ids", log.DeleteLoginLog)
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
