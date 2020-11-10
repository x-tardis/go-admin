package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/system"
)

func Menu(v1 gin.IRouter) {
	ctl := new(system.Menu)
	v1.GET("/menuTree/role", ctl.GetMenuTreeWithRoleName)
	r := v1.Group("/menus")
	{
		r.GET("", ctl.QueryTree)
		r.GET("/:id", ctl.Get)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:id", ctl.Delete)
	}
}

func PubMenu(v1 gin.IRouter) {
	v1.GET("/menuTree/option", new(system.Menu).GetMenuTreeOption)
}
