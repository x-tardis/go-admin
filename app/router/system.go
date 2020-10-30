package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/authj"
	"github.com/thinkgos/gin-middlewares/expvar"

	"github.com/x-tardis/go-admin/app/apis/system"
	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/routers"
)

func RegisterSystem(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	// public
	engine.GET("/", system.HelloWorld)
	engine.POST("/login", authMiddleware.LoginHandler)
	engine.GET("/refresh_token", authMiddleware.RefreshHandler) // Refresh time can be longer than token timeout
	engine.GET("/debug/vars", expvar.Handler())

	// 静态文件
	StaticFile(engine)
	// swagger
	routers.Swagger(engine)

	v1 := engine.Group("/api/v1")
	{ // 无需认证
		routers.PubSystem(v1)
		routers.PubPublic(v1)
		routers.PubBase(v1)
		routers.PubDB(v1)
		routers.PubSysTable(v1)

		// 需要认证
		v1.Use(
			authMiddleware.MiddlewareFunc(),
			OperLog(),
			authj.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.CasbinSubject),
		)
		v1.POST("/logout", authMiddleware.LogoutHandler)
		routers.System(v1)
		routers.Base(v1)
		routers.Dept(v1)
		routers.Dict(v1)
		routers.User(v1)
		routers.Role(v1)
		routers.Config(v1)
		routers.Post(v1)
		routers.Menu(v1)
		routers.LoginLog(v1)
		routers.OperLog(v1)
	}
}
