package router

import (
	"mime"

	"github.com/gin-gonic/gin"
)

func StaticFile(r gin.IRouter) {
	mime.AddExtensionType(".js", "application/javascript")
	r.Static("/static", "./static")
	r.Static("/form-generator", "./static/form-generator")
}
