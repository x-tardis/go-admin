package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/system"
)

func Role(v1 gin.IRouter) {
	ctl := new(system.Role)
	r := v1.Group("/roles")
	{
		r.GET("", ctl.QueryPage)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:ids", ctl.BatchDelete)
	}
}