package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/public"
)

func PubPublic(v1 gin.IRouter) {
	p := v1.Group("/public")
	{
		p.POST("/uploadFile", public.UploadFile)
	}
}
