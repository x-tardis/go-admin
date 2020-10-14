package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/system"
)

func Role(v1 gin.IRouter) {
	r := v1.Group("/role")
	{
		r.GET("/:roleId", system.GetRole)
		r.POST("", system.InsertRole)
		r.PUT("", system.UpdateRole)
		r.DELETE("/:roleId", system.DeleteRole)
	}
}
