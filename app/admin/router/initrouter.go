package router

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/admin/middleware"
	"github.com/x-tardis/go-admin/app/admin/middleware/handler"
	"github.com/x-tardis/go-admin/common/global"
	"github.com/x-tardis/go-admin/logger"
	"github.com/x-tardis/go-admin/pkg/deployed"
	_ "github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/tools"
)

func InitRouter() {
	var r *gin.Engine
	h := global.Cfg.GetEngine()
	if h == nil {
		h = gin.New()
		global.Cfg.SetEngine(h)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		logger.Fatal("not support other engine")
		os.Exit(-1)
	}
	if deployed.SslConfig.Enable {
		r.Use(handler.TlsHandler())
	}
	r.Use(middleware.WithContextDb(middleware.GetGormFromConfig(global.Cfg)))
	middleware.InitMiddleware(r)
	// the jwt middleware
	var err error
	authMiddleware, err := middleware.AuthInit()
	tools.HasError(err, "JWT Init Error", 500)

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	// 注册业务路由
	// TODO: 这里可存放业务路由，里边并无实际路由只有演示代码
	InitExamplesRouter(r, authMiddleware)

	//return r
}
