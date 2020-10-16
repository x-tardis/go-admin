package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Config(v1 gin.IRouter) {
	v1.GET("/configKey/:configKey", system.GetConfigByConfigKey)
	r := v1.Group("/configs")
	{
		r.GET("", system.GetConfigList)
		r.GET("/:id", system.GetConfig)
		r.POST("", system.InsertConfig)
		r.PUT("", system.UpdateConfig)
		r.DELETE("/:ids", system.DeleteConfig)
	}
}
