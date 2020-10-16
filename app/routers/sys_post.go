package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Post(v1 gin.IRouter) {
	r := v1.Group("/posts")
	{
		r.GET("", system.GetPostList)
		r.GET("/:id", system.GetPost)
		r.POST("", system.InsertPost)
		r.PUT("", system.UpdatePost)
		r.DELETE("/:ids", system.DeletePost)
	}
}
