package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/monitor"
	"github.com/x-tardis/go-admin/app/apis/system"
)

func PubSysSetting(v1 gin.IRouter) {
	r := v1.Group("/setting")
	{
		r.GET("", system.GetSetting)
		r.POST("", system.CreateSetting)
		r.GET("/serverInfo", monitor.ServerInfo)
	}
}
