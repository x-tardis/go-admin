package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

// Dict 字典类型与数据路由
func Dict(v1 gin.IRouter) {
	r := v1.Group("/dict")
	{
		ctlType := new(system.DictType)

		r.GET("/typeselect", ctlType.GetOptionSelect)
		tr := r.Group("/type")
		{
			tr.GET("", ctlType.QueryPage)
			tr.GET("/:dictId", ctlType.Get)
			tr.POST("", ctlType.Create)
			tr.PUT("", ctlType.Update)
			tr.DELETE("/:dictIds", ctlType.BatchDelete)
		}

		ctlData := new(system.DictData)
		dr := r.Group("/data")
		{
			dr.GET("", ctlData.QueryPage)
			dr.GET("/:dictId", ctlData.Get)
			dr.POST("", ctlData.Create)
			dr.PUT("", ctlData.Update)
			dr.DELETE("/:dictIds", ctlData.BatchDelete)
		}
	}
}
