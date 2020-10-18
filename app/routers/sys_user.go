package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func User(v1 gin.IRouter) {
	ctl := new(system.User)
	r := v1.Group("/users")
	{
		r.GET("", ctl.QueryPage)
		r.GET("/:id", ctl.Get)
		r.GET("/", ctl.GetInit)
		r.POST("", ctl.Create)
		r.PUT("", ctl.Update)
		r.DELETE("/:ids", ctl.BatchDelete)
	}
	rx := v1.Group("/user")
	{
		rx.GET("/profile", ctl.GetProfile)
		rx.POST("/avatar", ctl.UploadAvatar)
		rx.PUT("/pwd", ctl.UpdatePassword)
	}
}
