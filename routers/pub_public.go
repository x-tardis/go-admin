package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/public"
)

func PubPublic(v1 gin.IRouter) {
	r := v1.Group("/public")
	{
		r.POST("/uploadFile", public.UploadFile)
	}
}
