package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/monitor"
	"github.com/x-tardis/go-admin/app/apis/system"
)

func PubSysSetting(v1 gin.IRouter) {
	r := v1.Group("/system")
	{
		r.GET("/ping", system.Ping)
		r.GET("/setting", system.GetSetting)
		r.POST("/setting", system.CreateSetting)
		r.GET("/info", monitor.SystemInfo)
	}
}
