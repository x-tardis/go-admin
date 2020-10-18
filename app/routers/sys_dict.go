package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/dict"
)

func Dict(v1 *gin.RouterGroup) {
	ctl := new(dict.DictData)
	r := v1.Group("/dict")
	{
		r.GET("/data", ctl.QueryPage)
		r.GET("/data/:dictCode", ctl.Get)
		r.POST("/data", ctl.Create)
		r.PUT("/data", ctl.Update)
		r.DELETE("/data/:dictCode", ctl.BatchDelete)

		r.GET("/typelist", dict.GetDictTypeList)
		r.GET("/type/:dictId", dict.GetDictType)
		r.POST("/type", dict.InsertDictType)
		r.PUT("/type", dict.UpdateDictType)
		r.DELETE("/type/:dictId", dict.DeleteDictType)

		r.GET("/typeoptionselect", dict.GetDictTypeOptionSelect)
	}
}
