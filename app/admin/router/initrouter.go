package router

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/admin/middleware"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/logger"
	middleware2 "github.com/x-tardis/go-admin/pkg/middleware"
	"github.com/x-tardis/go-admin/tools"
)

func InitRouter() {
	var r *gin.Engine
	h := deployed.Cfg.GetEngine()
	if h == nil {
		h = gin.New()
		deployed.Cfg.SetEngine(h)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		logger.Fatal("not support other engine")
		os.Exit(-1)
	}
	if deployed.SslConfig.Enable {
		r.Use(middleware2.Tls(deployed.SslConfig.Domain))
	}
	r.Use(middleware2.WithContextDb(middleware2.GetGormFromConfig(deployed.Cfg)))
	r.Use(middleware.LoggerToFile(), // 日志处理
		middleware2.Recovery(), // 自定义错误处理
		middleware2.NoCache(),  // NoCache is a middleware function that appends headers
		middleware2.Cors(),     // 跨域处理
		middleware2.Secure(),   // Secure is a middleware function that appends security
	)
	// the jwt middleware
	var err error
	authMiddleware, err := middleware.AuthInit(deployed.ApplicationConfig.JwtSecret)
	tools.HasError(err, "JWT Init Error", 500)

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	// 注册业务路由
	// TODO: 这里可存放业务路由，里边并无实际路由只有演示代码
	InitExamplesRouter(r, authMiddleware)

	//return r
}
