package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func RoleDept(v1 gin.IRouter) {
	ctl := new(system.RoleDept)
	v1.GET("/roleDeptTree/option/:roleId", ctl.GetDeptTreeOptionRole)
}
