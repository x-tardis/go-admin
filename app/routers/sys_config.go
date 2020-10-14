package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Config(v1 gin.IRouter) {
	v1.GET("/configKey/:configKey", system.GetConfigByConfigKey)
	v1.GET("/configList", system.GetConfigList)
	r := v1.Group("/config")
	{
		r.GET("/:configId", system.GetConfig)
		r.POST("", system.InsertConfig)
		r.PUT("", system.UpdateConfig)
		r.DELETE("/:configId", system.DeleteConfig)
	}
}
