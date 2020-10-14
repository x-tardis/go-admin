package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Post(v1 gin.IRouter) {
	r := v1.Group("/post")
	{
		r.GET("/:postId", system.GetPost)
		r.POST("", system.InsertPost)
		r.PUT("", system.UpdatePost)
		r.DELETE("/:postId", system.DeletePost)
	}
}
