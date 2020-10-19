package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/dict"
	"github.com/x-tardis/go-admin/app/apis/system"
	"github.com/x-tardis/go-admin/app/apis/tools"
)

func PubBase(v1 gin.IRouter) {
	v1.GET("/captcha", system.GetCaptcha)
	v1.GET("/menuTreeselect", system.GetMenuTreeselect)
	v1.GET("/dict/databytype/:dictType", new(dict.DictData).GetWithType)
	r := v1.Group("/gen")
	{
		r.GET("/preview/:tableId", tools.Preview)
		r.GET("/toproject/:tableId", tools.GenCodeV3)
		r.GET("/todb/:tableId", tools.GenMenuAndApi)
		r.GET("/tabletree", tools.GetSysTablesTree)
	}
}
