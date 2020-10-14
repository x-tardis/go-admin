package routers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/log"
	"github.com/x-tardis/go-admin/app/apis/system"
)

func Page(v1 gin.IRouter) {
	v1.GET("/deptList", system.GetDeptList)
	v1.GET("/deptTree", system.GetDeptTree)
	v1.GET("/sysUserList", system.GetSysUserList)
	v1.GET("/rolelist", system.GetRoleList)
	v1.GET("/configList", system.GetConfigList)
	v1.GET("/postlist", system.GetPostList)
	v1.GET("/menulist", system.GetMenuList)
	v1.GET("/loginloglist", log.GetLoginLogList)
}

func Base(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1.GET("/getinfo", system.GetInfo)
	v1.GET("/menurole", system.GetMenuRole)
	v1.PUT("/roledatascope", system.UpdateRoleDataScope)
	v1.GET("/roleMenuTreeselect/:roleId", system.GetMenuTreeRoleselect)
	v1.GET("/roleDeptTreeselect/:roleId", system.GetDeptTreeRoleselect)

	v1.POST("/logout", authMiddleware.LogoutHandler)
	v1.GET("/menuids", system.GetMenuIDS)

	v1.GET("/operloglist", log.GetOperLogList)
	v1.GET("/configKey/:configKey", system.GetConfigByConfigKey)
}
