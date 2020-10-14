package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/routers"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/middleware"
)

var routerNoCheckRole []func(r gin.IRouter)
var routerCheckRole []func(r gin.IRouter)

func init() {
	routerNoCheckRole = append(routerNoCheckRole,
		routers.PubSysFileInfo,
		routers.PubSysFileDir,
	)
	routerCheckRole = append(routerCheckRole,
		routers.SysJobRouter,
		routers.SysContent,
		routers.SysCategory,
	)
}

// 路由示例
func RegisterBusiness(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.Engine {
	v1 := r.Group("/api/v1")
	{ // 无需认证的路由
		for _, f := range routerNoCheckRole {
			f(v1)
		}
	}
	{ // 需要认证的路由
		v1.Use(authMiddleware.MiddlewareFunc(), middleware.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.RoleKey))
		for _, f := range routerCheckRole {
			f(v1)
		}
	}

	return r
}
