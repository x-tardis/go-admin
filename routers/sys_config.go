package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/system"
)

func Config(v1 gin.IRouter) {
	ctl := new(system.Config)
	v1.GET("/configKey/:key", ctl.GetWithKey)
	r := v1.Group("/configs")
	{
		r.GET("", ctl.QueryPage)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:ids", ctl.BatchDelete)
	}
}
