package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/authj"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/routers"
)

// 业务无授权,无RBAC角色控制
var businessNoAuthRbac []func(r gin.IRouter)

// 业务有授权,有RBAC角色控制
var businessAuthRbac []func(r gin.IRouter)

func init() {
	businessNoAuthRbac = append(businessNoAuthRbac,
		routers.PubSysFileInfo,
		routers.PubSysFileDir,
	)
	businessAuthRbac = append(businessAuthRbac,
		routers.SysJobRouter,
		routers.SysContent,
		routers.SysCategory,
	)
}

// 路由示例
func RegisterBusiness(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := engine.Group("/api/v1")
	{
		// 无需认证的路由
		for _, f := range businessNoAuthRbac {
			f(v1)
		}
		// 需要认证的路由
		v1.Use(
			authMiddleware.MiddlewareFunc(),
			OperLog(),
			authj.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.CasbinSubject),
		)
		for _, f := range businessAuthRbac {
			f(v1)
		}
	}
}
