package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/system"
)

func Role(v1 gin.IRouter) {
	ctl := new(system.Role)
	v1.PUT("/roledatascope", ctl.UpdateDataScope)
	r := v1.Group("/roles")
	{
		r.GET("", ctl.QueryPage)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:ids", ctl.BatchDelete)
		r.PATCH("/enable/:id", ctl.Enable)
		r.PATCH("/disable/:id", ctl.Disable)
	}
}
