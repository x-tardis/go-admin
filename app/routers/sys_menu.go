package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Menu(v1 gin.IRouter) {
	ctl := new(system.Menu)
	r := v1.Group("/menus")
	{
		r.GET("", ctl.QueryTree)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:id", ctl.Delete)
	}
}
