package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/tools"
)

func PubSysTable(v1 gin.IRouter) {
	r := v1.Group("/sys/tables")
	{
		r.GET("/page", tools.GetSysTableList)
		tablesInfo := r.Group("/info")
		{
			tablesInfo.GET("", tools.GetSysTablesInfo)
			tablesInfo.GET("/:id", tools.GetSysTables)
			tablesInfo.POST("", tools.InsertSysTable)
			tablesInfo.PUT("", tools.UpdateSysTable)
			tablesInfo.DELETE("/:ids", tools.DeleteSysTables)
		}
	}
}
