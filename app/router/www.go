package router

import (
	"mime"

	"github.com/gin-gonic/gin"

	_ "github.com/x-tardis/go-admin/docs"
)

func sysStaticFileRouter(r *gin.RouterGroup) {
	mime.AddExtensionType(".js", "application/javascript")
	r.Static("/static", "./static")
	r.Static("/form-generator", "./static/form-generator")
}
