package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/system"
)

func Dept(v1 gin.IRouter) {
	ctl := new(system.Dept)
	v1.GET("/deptTree", ctl.QueryTree)
	r := v1.Group("/depts")
	{
		r.GET("", ctl.QueryPage)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:id", ctl.Delete)
	}
}
