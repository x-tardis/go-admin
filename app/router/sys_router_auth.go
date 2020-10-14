package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	log2 "github.com/x-tardis/go-admin/app/apis/log"
	"github.com/x-tardis/go-admin/app/apis/system"
	"github.com/x-tardis/go-admin/app/apis/system/dict"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/middleware"
	"github.com/x-tardis/go-admin/pkg/ws"
)

func sysCheckRoleRouterInit(r gin.IRouter, authMiddleware *jwt.GinJWTMiddleware) {
	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler) // Refresh time can be longer than token timeout

	wsGroup := r.Group("")
	wsGroup.Use(authMiddleware.MiddlewareFunc())
	{
		wsGroup.GET("/ws/:id/:channel", ws.WebsocketManager.WsClient)
		wsGroup.GET("/wslogout/:id/:channel", ws.WebsocketManager.UnWsClient)
	}

	v1 := r.Group("/api/v1")
	v1.Use(authMiddleware.MiddlewareFunc(), middleware.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.RoleKey))
	{
		registerPageRouter(v1)
		registerBaseRouter(v1, authMiddleware)
		registerDeptRouter(v1)
		registerDictRouter(v1)
		registerSysUserRouter(v1)
		registerRoleRouter(v1)
		registerConfigRouter(v1)
		registerUserCenterRouter(v1)
		registerPostRouter(v1)
		registerMenuRouter(v1)
		registerLoginLogRouter(v1)
		registerOperLogRouter(v1)
	}
}

func registerPageRouter(v1 gin.IRouter) {
	v1.GET("/deptList", system.GetDeptList)
	v1.GET("/deptTree", system.GetDeptTree)
	v1.GET("/sysUserList", system.GetSysUserList)
	v1.GET("/rolelist", system.GetRoleList)
	v1.GET("/configList", system.GetConfigList)
	v1.GET("/postlist", system.GetPostList)
	v1.GET("/menulist", system.GetMenuList)
	v1.GET("/loginloglist", log2.GetLoginLogList)
}

func registerBaseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1.GET("/getinfo", system.GetInfo)
	v1.GET("/menurole", system.GetMenuRole)
	v1.PUT("/roledatascope", system.UpdateRoleDataScope)
	v1.GET("/roleMenuTreeselect/:roleId", system.GetMenuTreeRoleselect)
	v1.GET("/roleDeptTreeselect/:roleId", system.GetDeptTreeRoleselect)

	v1.POST("/logout", authMiddleware.LogoutHandler)
	v1.GET("/menuids", system.GetMenuIDS)

	v1.GET("/operloglist", log2.GetOperLogList)
	v1.GET("/configKey/:configKey", system.GetConfigByConfigKey)
}

func registerDeptRouter(v1 *gin.RouterGroup) {
	dept := v1.Group("/dept")
	{
		dept.GET("/:deptId", system.GetDept)
		dept.POST("", system.InsertDept)
		dept.PUT("", system.UpdateDept)
		dept.DELETE("/:id", system.DeleteDept)
	}
}

func registerDictRouter(v1 *gin.RouterGroup) {
	dicts := v1.Group("/dict")
	{
		dicts.GET("/datalist", dict.GetDictDataList)
		dicts.GET("/typelist", dict.GetDictTypeList)
		dicts.GET("/typeoptionselect", dict.GetDictTypeOptionSelect)

		dicts.GET("/data/:dictCode", dict.GetDictData)
		dicts.POST("/data", dict.InsertDictData)
		dicts.PUT("/data/", dict.UpdateDictData)
		dicts.DELETE("/data/:dictCode", dict.DeleteDictData)

		dicts.GET("/type/:dictId", dict.GetDictType)
		dicts.POST("/type", dict.InsertDictType)
		dicts.PUT("/type", dict.UpdateDictType)
		dicts.DELETE("/type/:dictId", dict.DeleteDictType)
	}
}

func registerSysUserRouter(v1 gin.IRouter) {
	sysuser := v1.Group("/sysUser")
	{
		sysuser.GET("/:userId", system.GetSysUser)
		sysuser.GET("/", system.GetSysUserInit)
		sysuser.POST("", system.InsertSysUser)
		sysuser.PUT("", system.UpdateSysUser)
		sysuser.DELETE("/:userId", system.DeleteSysUser)
	}
}

func registerRoleRouter(v1 gin.IRouter) {
	role := v1.Group("/role")
	{
		role.GET("/:roleId", system.GetRole)
		role.POST("", system.InsertRole)
		role.PUT("", system.UpdateRole)
		role.DELETE("/:roleId", system.DeleteRole)
	}
}

func registerConfigRouter(v1 gin.IRouter) {
	config := v1.Group("/config")
	{
		config.GET("/:configId", system.GetConfig)
		config.POST("", system.InsertConfig)
		config.PUT("", system.UpdateConfig)
		config.DELETE("/:configId", system.DeleteConfig)
	}
}

func registerUserCenterRouter(v1 gin.IRouter) {
	user := v1.Group("/user")
	{
		user.GET("/profile", system.GetSysUserProfile)
		user.POST("/avatar", system.InsetSysUserAvatar)
		user.PUT("/pwd", system.SysUserUpdatePwd)
	}
}

func registerPostRouter(v1 gin.IRouter) {
	post := v1.Group("/post")
	{
		post.GET("/:postId", system.GetPost)
		post.POST("", system.InsertPost)
		post.PUT("", system.UpdatePost)
		post.DELETE("/:postId", system.DeletePost)
	}
}

func registerMenuRouter(v1 gin.IRouter) {
	menu := v1.Group("/menu")
	{
		menu.GET("/:id", system.GetMenu)
		menu.POST("", system.InsertMenu)
		menu.PUT("", system.UpdateMenu)
		menu.DELETE("/:id", system.DeleteMenu)
	}
}
func registerLoginLogRouter(v1 gin.IRouter) {
	loginlog := v1.Group("/loginlog")
	{
		loginlog.GET("/:infoId", log2.GetLoginLog)
		loginlog.POST("", log2.InsertLoginLog)
		loginlog.PUT("", log2.UpdateLoginLog)
		loginlog.DELETE("/:infoId", log2.DeleteLoginLog)
	}
}

func registerOperLogRouter(v1 gin.IRouter) {
	operlog := v1.Group("/operlog")
	{
		operlog.GET("/:operId", log2.GetOperLog)
		operlog.DELETE("/:operId", log2.DeleteOperLog)
	}
}
