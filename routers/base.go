package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Base(v1 gin.IRouter) {
	v1.GET("/menurole", system.GetMenuTreeWithRoleName)
	v1.GET("/menuids", system.GetMenuIDS)
	v1.PUT("/roledatascope", system.UpdateRoleDataScope)
	v1.GET("/roleMenuTreeoption/:roleId", system.GetMenuTreeRoleOption)
	v1.GET("/roleDeptTreeoption/:roleId", system.GetDeptTreeRoleOption)
}
