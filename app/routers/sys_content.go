package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/syscontent"
)

// 需认证的路由代码
func SysContent(v1 gin.IRouter) {
	r := v1.Group("/syscontents")
	{
		r.GET("", syscontent.GetSysContentList)
		r.GET("/:id", syscontent.GetSysContent)
		r.POST("", syscontent.InsertSysContent)
		r.PUT("", syscontent.UpdateSysContent)
		r.DELETE("/:id", syscontent.DeleteSysContent)
	}
}
