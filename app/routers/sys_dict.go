package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system/dict"
)

func Dict(v1 *gin.RouterGroup) {
	r := v1.Group("/dict")
	{
		r.GET("/datalist", dict.GetDictDataList)
		r.GET("/typelist", dict.GetDictTypeList)
		r.GET("/typeoptionselect", dict.GetDictTypeOptionSelect)

		r.GET("/data/:dictCode", dict.GetDictData)
		r.POST("/data", dict.InsertDictData)
		r.PUT("/data/", dict.UpdateDictData)
		r.DELETE("/data/:dictCode", dict.DeleteDictData)

		r.GET("/type/:dictId", dict.GetDictType)
		r.POST("/type", dict.InsertDictType)
		r.PUT("/type", dict.UpdateDictType)
		r.DELETE("/type/:dictId", dict.DeleteDictType)
	}
}
