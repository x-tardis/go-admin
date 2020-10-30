package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/tools"
)

func PubGen(v1 gin.IRouter) {
	r := v1.Group("/gen")
	{
		r.GET("/preview/:tableId", tools.Preview)
		r.GET("/toproject/:tableId", tools.GenCodeV3)
		r.GET("/todb/:tableId", tools.GenMenuAndApi)
	}
}
