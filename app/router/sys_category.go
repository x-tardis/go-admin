package router

import (
	"github.com/gin-gonic/gin"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/x-tardis/go-admin/app/apis/syscategory"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/middleware"
)

// 需认证的路由代码
func registerSysCategoryRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/syscategory").
		Use(authMiddleware.MiddlewareFunc(), middleware.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.RoleKey))
	{
		r.GET("/:id", syscategory.GetSysCategory)
		r.POST("", syscategory.InsertSysCategory)
		r.PUT("", syscategory.UpdateSysCategory)
		r.DELETE("/:id", syscategory.DeleteSysCategory)
	}

	l := v1.Group("").
		Use(authMiddleware.MiddlewareFunc(), middleware.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.RoleKey))
	{
		l.GET("/syscategoryList", syscategory.GetSysCategoryList)
	}

}
