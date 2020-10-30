package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Base(v1 gin.IRouter) {
	v1.PUT("/roledatascope", system.UpdateRoleDataScope)
}
