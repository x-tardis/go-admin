package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
	"github.com/x-tardis/go-admin/app/apis/system/dict"
	"github.com/x-tardis/go-admin/app/apis/tools"
)

func PubBase(v1 gin.IRouter) {
	v1.GET("/getCaptcha", system.GenerateCaptchaHandler)
	v1.GET("/gen/preview/:tableId", tools.Preview)
	v1.GET("/gen/toproject/:tableId", tools.GenCodeV3)
	v1.GET("/gen/todb/:tableId", tools.GenMenuAndApi)
	v1.GET("/gen/tabletree", tools.GetSysTablesTree)
	v1.GET("/menuTreeselect", system.GetMenuTreeelect)
	v1.GET("/dict/databytype/:dictType", dict.GetDictDataByDictType)
}
