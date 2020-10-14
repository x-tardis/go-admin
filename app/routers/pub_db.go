package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/tools"
)

func PubDB(v1 gin.IRouter) {
	r := v1.Group("/db")
	{
		r.GET("/tables/page", tools.GetDBTableList)
		r.GET("/columns/page", tools.GetDBColumnList)
	}
}
