package router

import (
	"context"
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/didip/tollbooth/v5"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thinkgos/gin-middlewares/authj"
	"github.com/thinkgos/gin-middlewares/expvar"
	"github.com/thinkgos/gin-middlewares/gzap"
	"github.com/thinkgos/gin-middlewares/nocache"
	"github.com/thinkgos/gin-middlewares/pprof"
	"github.com/thinkgos/gin-middlewares/ratelimiter"
	"github.com/thinkgos/gin-middlewares/requestid"

	"github.com/x-tardis/go-admin/apis/system"
	"github.com/x-tardis/go-admin/apis/ws"
	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/middleware"
	"github.com/x-tardis/go-admin/pkg/xxfield"
	"github.com/x-tardis/go-admin/routers"
)

func InitRouter() *gin.Engine {
	engine := gin.New()

	if deployed.SslConfig.Enable {
		engine.Use(middleware.Tls(deployed.SslConfig.Domain))
	}

	engine.Use(
		requestid.RequestID(), // request id
		gzap.Logger(deployed.RequestLogger, gzap.WithCustomFields(xxfield.RequestId, xxfield.Error)),            // logger
		gzap.Recovery(deployed.RequestLogger, !deployed.IsModeProd(), gzap.WithCustomFields(xxfield.RequestId)), // recover, 仅开发时开启stack
		nocache.NoCache(),              // NoCache is a middleware function that appends headers
		cors.New(deployed.ViperCors()), // 跨域处理
		ratelimiter.RateLimit(tollbooth.NewLimiter(deployed.ViperLimiter(), nil)), // 限速器
		middleware.Secure(),                // Secure is a middleware function that appends security
		middleware.WriteResponseHeaderID(), // 写应答头部写request id
	)

	// the jwt middleware
	authMiddleware, err := system.NewJWTAuth(deployed.ViperJwt())
	if err != nil {
		log.Fatalf("jwt initialize failed, %+v", err)
	}

	RegisterPubic(engine, authMiddleware)  // 注册公共开放接口
	RegisterSystem(engine, authMiddleware) // 注册系统路由
	RegisterWs(engine, authMiddleware)     // 注册ws
	custom(engine, authMiddleware)
	return engine
}

func RegisterPubic(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	// public
	engine.GET("/", system.HelloWorld)
	engine.POST("/login", authMiddleware.LoginHandler)
	engine.GET("/refresh_token", authMiddleware.RefreshHandler) // Refresh time can be longer than token timeout
	// metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// debug
	engine.GET("/debug/vars", expvar.Handler())
	pprof.Router(engine)

	// 静态文件
	StaticFile(engine)
	// swagger
	routers.Swagger(engine)
}

func RegisterSystem(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := engine.Group("/api/v1")
	{
		// 无需认证
		routers.PubBase(v1)
		routers.PubSystem(v1)
		routers.PubPublic(v1)
		routers.PubMenu(v1)
		routers.PubGen(v1)
		routers.PubDB(v1)
		routers.PubSysTable(v1)
		// business public
		businessMu.RLock()
		for _, f := range businessPublic {
			f(v1)
		}
		businessMu.RUnlock()
	}
	{
		// 需要认证
		v1.Use(
			authMiddleware.MiddlewareFunc(),
			OperLog(),
			authj.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.CasbinSubject),
		)
		v1.POST("/logout", authMiddleware.LogoutHandler)
		routers.System(v1)
		routers.Config(v1)
		routers.Dict(v1)
		routers.LoginLog(v1)
		routers.OperLog(v1)
		routers.User(v1)
		routers.Menu(v1)
		routers.Role(v1)
		routers.RoleDept(v1)
		routers.RoleMenu(v1)
		routers.Dept(v1)
		routers.Post(v1)
		// business
		businessMu.RLock()
		for _, f := range businessAuthRbac {
			f(v1)
		}
		businessMu.RUnlock()
	}
}

func RegisterWs(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	go ws.WsHub.Run(context.Background())
	go ws.FileMonitor(context.Background(), "temp/job.log", ws.JobGroup, ws.SendGroup)
	// 需要认证
	wsGroup := engine.Group("/ws")
	{
		wsGroup.Use(authMiddleware.MiddlewareFunc(), OperLog())
		routers.WsJob(wsGroup)
	}
}
