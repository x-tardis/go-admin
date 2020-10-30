package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func PubBase(v1 gin.IRouter) {
	v1.GET("/captcha", system.GetCaptcha)
	// v1.GET("/dict/databytype", nil)
	v1.GET("/dict/databytype/:dictType", new(system.DictData).GetWithType)
}
