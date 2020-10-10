package router

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/admin/apis/syscategory"
	"github.com/x-tardis/go-admin/pkg/deployed"
	jwt "github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/middleware"
)

// 需认证的路由代码
func registerSysCategoryRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r := v1.Group("/syscategory").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(deployed.CasbinEnforcer))
	{
		r.GET("/:id", syscategory.GetSysCategory)
		r.POST("", syscategory.InsertSysCategory)
		r.PUT("", syscategory.UpdateSysCategory)
		r.DELETE("/:id", syscategory.DeleteSysCategory)
	}

	l := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(deployed.CasbinEnforcer))
	{
		l.GET("/syscategoryList", syscategory.GetSysCategoryList)
	}

}
