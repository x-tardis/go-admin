package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Menu(v1 gin.IRouter) {
	r := v1.Group("/menus")
	{
		r.GET("", system.GetMenuList)
		r.GET("/:id", system.GetMenu)
		r.POST("", system.InsertMenu)
		r.PUT("", system.UpdateMenu)
		r.DELETE("/:id", system.DeleteMenu)
	}
}
