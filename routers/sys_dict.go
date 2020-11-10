package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/system"
)

// Dict 字典类型与数据路由
func Dict(v1 gin.IRouter) {
	r := v1.Group("/dict")
	{
		ctlType := new(system.DictType)
		r.GET("/typeoption", ctlType.GetOption)
		tr := r.Group("/type")
		{
			tr.GET("", ctlType.QueryPage)
			tr.GET("/:id", ctlType.Get)
			tr.POST("", ctlType.Create)
			tr.PUT("", ctlType.Update)
			tr.DELETE("/:ids", ctlType.BatchDelete)
		}

		ctlData := new(system.DictData)
		dr := r.Group("/data")
		{
			dr.GET("", ctlData.QueryPage)
			dr.GET("/:id", ctlData.Get)
			dr.POST("", ctlData.Create)
			dr.PUT("", ctlData.Update)
			dr.DELETE("/:ids", ctlData.BatchDelete)
		}
	}
}
