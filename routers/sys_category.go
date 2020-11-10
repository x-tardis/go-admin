package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/syscategory"
)

// 需认证的路由代码
func SysCategory(v1 gin.IRouter) {
	ctl := new(syscategory.Category)
	r := v1.Group("/categories")
	{
		r.GET("", ctl.QueryPage)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:id", ctl.BatchDelete)
	}
}
