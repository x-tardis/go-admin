package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/gzap"
	"github.com/thinkgos/gin-middlewares/requestid"

	"github.com/x-tardis/go-admin/app/apis/system"
	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/middleware"
)

func InitRouter() *gin.Engine {
	engine := gin.New()

	if deployed.SslConfig.Enable {
		engine.Use(middleware.Tls(deployed.SslConfig.Domain))
	}

	engine.Use(
		requestid.RequestID(),                                                // request id
		gzap.Logger(deployed.RequestLogger.Desugar()),                        // logger
		gzap.Recovery(izap.Logger, deployed.AppConfig.Mode == infra.ModeDev), // recover, 仅开发时开启stack
		middleware.NoCache(),                                                 // NoCache is a middleware function that appends headers
		cors.New(*deployed.CorsConfig),                                       // 跨域处理
		middleware.Secure(),                                                  // Secure is a middleware function that appends security
	)

	// the jwt middleware
	authMiddleware, err := system.NewJWTAuth(deployed.JwtConfig)
	if err != nil {
		panic("jwt int failed")
	}

	RegisterSystem(engine, authMiddleware)   // 注册系统路由
	RegisterWs(engine, authMiddleware)       // 注册ws
	RegisterBusiness(engine, authMiddleware) // 注册业务路由
	return engine
}
