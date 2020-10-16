package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/syscategory"
)

// 需认证的路由代码
func SysCategory(v1 gin.IRouter) {
	r := v1.Group("/syscategory")
	{
		r.GET("", syscategory.GetSysCategoryList)
		r.GET("/:id", syscategory.GetSysCategory)
		r.POST("", syscategory.InsertSysCategory)
		r.PUT("", syscategory.UpdateSysCategory)
		r.DELETE("/:id", syscategory.DeleteSysCategory)
	}
}
