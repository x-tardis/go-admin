package routers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/system"
)

func Base(v1 gin.IRouter, authMiddleware *jwt.GinJWTMiddleware) {
	v1.POST("/logout", authMiddleware.LogoutHandler)
	v1.GET("/menurole", system.GetMenuTreeWithRoleName)
	v1.GET("/menuids", system.GetMenuIDS)
	v1.PUT("/roledatascope", system.UpdateRoleDataScope)
	v1.GET("/roleMenuTreeoption/:roleId", system.GetMenuTreeRoleOption)
	v1.GET("/roleDeptTreeoption/:roleId", system.GetDeptTreeRoleOption)
}
