package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/monitor"
	"github.com/x-tardis/go-admin/app/apis/system"
)

func PubSystem(v1 gin.IRouter) {
	r := v1.Group("/system")
	{
		r.GET("/ping", system.Ping)
		r.GET("/info", monitor.SystemInfo)

		// system setting
		ctl := new(system.SysSetting)
		rx := r.Group("/setting")
		{
			rx.GET("", ctl.Get)
			rx.PUT("", ctl.Update)
		}
	}
}
