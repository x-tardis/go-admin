package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/x-tardis/go-admin/app/apis/system"
)

func RoleMenu(v1 gin.IRouter) {
	ctl := new(system.RoleMenu)
	v1.GET("/roleMenuTree/option/:roleId", ctl.GetMenuTreeOptionRole)
}
