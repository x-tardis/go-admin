package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Dict(v1 *gin.RouterGroup) {
	ctlType := new(system.DictType)
	ctlData := new(system.DictData)
	r := v1.Group("/dict")
	{
		r.GET("/data", ctlData.QueryPage)
		r.GET("/data/:dictCode", ctlData.Get)
		r.POST("/data", ctlData.Create)
		r.PUT("/data", ctlData.Update)
		r.DELETE("/data/:dictCode", ctlData.BatchDelete)

		r.GET("/type", ctlType.QueryPage)
		r.GET("/type/:dictId", ctlType.Get)
		r.POST("/type", ctlType.Create)
		r.PUT("/type", ctlType.Update)
		r.DELETE("/type/:dictId", ctlType.BatchDelete)
		r.GET("/typeoptionselect", ctlType.GetOptionSelect)
	}
}
