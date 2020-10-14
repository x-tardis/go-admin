package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/system"
)

func Dept(v1 gin.IRouter) {
	v1.GET("/deptList", system.GetDeptList)
	v1.GET("/deptTree", system.GetDeptTree)
	r := v1.Group("/dept")
	{
		r.GET("/:deptId", system.GetDept)
		r.POST("", system.InsertDept)
		r.PUT("", system.UpdateDept)
		r.DELETE("/:id", system.DeleteDept)
	}
}
