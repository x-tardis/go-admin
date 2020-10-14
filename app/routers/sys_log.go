package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/log"
)

func LoginLog(v1 gin.IRouter) {
	r := v1.Group("/loginlog")
	{
		r.GET("/:infoId", log.GetLoginLog)
		r.POST("", log.InsertLoginLog)
		r.PUT("", log.UpdateLoginLog)
		r.DELETE("/:infoId", log.DeleteLoginLog)
	}
}
func OperLog(v1 gin.IRouter) {
	r := v1.Group("/operlog")
	{
		r.GET("/:operId", log.GetOperLog)
		r.DELETE("/:operId", log.DeleteOperLog)
	}
}
