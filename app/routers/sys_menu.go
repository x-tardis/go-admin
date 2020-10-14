package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Menu(v1 gin.IRouter) {
	v1.GET("/menulist", system.GetMenuList)
	r := v1.Group("/menu")
	{
		r.GET("/:id", system.GetMenu)
		r.POST("", system.InsertMenu)
		r.PUT("", system.UpdateMenu)
		r.DELETE("/:id", system.DeleteMenu)
	}
}
