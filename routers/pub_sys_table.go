package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/tools"
)

func PubSysTable(v1 gin.IRouter) {
	ctl := tools.Tables{}
	r := v1.Group("/sys/tables")
	{
		r.GET("/page", ctl.QueryPage)
		r.GET("/tree", ctl.QueryTree)
		tablesInfo := r.Group("/info")
		{
			tablesInfo.GET("", ctl.GetWithName)
			tablesInfo.GET("/:id", ctl.Get)
			tablesInfo.POST("", ctl.Create)
			tablesInfo.PUT("", ctl.Update)
			tablesInfo.DELETE("/:ids", ctl.DeleteSysTables)
		}
	}
}
