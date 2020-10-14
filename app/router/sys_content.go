package router

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/syscontent"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/x-tardis/go-admin/pkg/middleware"
)

// 需认证的路由代码
func registerSysContentRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/syscontent").
		Use(authMiddleware.MiddlewareFunc(), middleware.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.RoleKey))
	{
		r.GET("/:id", syscontent.GetSysContent)
		r.POST("", syscontent.InsertSysContent)
		r.PUT("", syscontent.UpdateSysContent)
		r.DELETE("/:id", syscontent.DeleteSysContent)
	}

	l := v1.Group("").
		Use(authMiddleware.MiddlewareFunc(), middleware.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.RoleKey))
	{
		l.GET("/syscontentList", syscontent.GetSysContentList)
	}
}
