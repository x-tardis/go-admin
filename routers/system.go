package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/system"
)

func PubSystem(v1 gin.IRouter) {
	r := v1.Group("/system")
	{
		r.GET("/ping", system.Ping)
		r.GET("/info", system.SystemInfo)
		r.GET("/setting", new(system.Setting).Get)
	}
}

func System(v1 gin.IRouter) {
	v1.PUT("/system/setting", new(system.Setting).Update)
}
