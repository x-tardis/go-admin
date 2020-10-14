package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/tools"
)

func SysTable(v1 gin.IRouter) {
	r := v1.Group("/sys/tables")
	{
		r.GET("/page", tools.GetSysTableList)
		tablesInfo := r.Group("/info")
		{
			tablesInfo.POST("", tools.InsertSysTable)
			tablesInfo.PUT("", tools.UpdateSysTable)
			tablesInfo.DELETE("/:tableId", tools.DeleteSysTables)
			tablesInfo.GET("/:tableId", tools.GetSysTables)
			tablesInfo.GET("", tools.GetSysTablesInfo)
		}
	}
}
