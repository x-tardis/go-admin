package router

import (
	"context"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"github.com/thinkgos/gin-middlewares/gzap"

	"github.com/x-tardis/go-admin/app/apis/auth"
	"github.com/x-tardis/go-admin/app/apis/system"
	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/app/routers"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/middleware"
	"github.com/x-tardis/go-admin/pkg/ws"
)

func InitRouter() *gin.Engine {
	engine := gin.New()

	if deployed.SslConfig.Enable {
		engine.Use(middleware.Tls(deployed.SslConfig.Domain))
	}

	engine.Use(
		middleware.WithContextDb(middleware.GetGormFromConfig(deployed.Cfg)),
		gzap.Logger(deployed.RequestLogger.Desugar()), // logger
		gzap.Recovery(izap.Logger, false),             // recover
		OperLog(),                                     // 操作日志写入数据库
		middleware.NoCache(),                          // NoCache is a middleware function that appends headers
		cors.New(*deployed.CorsConfig),                // 跨域处理
		middleware.Secure(),                           // Secure is a middleware function that appends security
	)
	// the jwt middleware
	authMiddleware, err := auth.NewJWTAuth(deployed.JwtConfig)
	if err != nil {
		panic("JWT Init Error")
	}

	// 注册系统路由
	RegisterSys(engine, authMiddleware)
	// 注册业务路由
	RegisterBusiness(engine, authMiddleware)

	return engine
}

func RegisterSys(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendAllService()

	engine.GET("/", system.HelloWorld)
	engine.POST("/login", authMiddleware.LoginHandler)
	engine.GET("/refresh_token", authMiddleware.RefreshHandler) // Refresh time can be longer than token timeout

	// 静态文件
	StaticFile(engine)
	// swagger
	routers.Swagger(engine)

	// 需要认证
	wsGroup := engine.Group("", authMiddleware.MiddlewareFunc())
	{
		wsGroup.GET("/ws/:id/:channel", ws.WebsocketManager.WsClient)
		wsGroup.GET("/wslogout/:id/:channel", ws.WebsocketManager.UnWsClient)
	}

	v1Group := engine.Group("/api/v1")
	{ // 无需认证
		routers.PubBase(v1Group)
		routers.PubDB(v1Group)
		routers.PubSysTable(v1Group)
		routers.PubPublic(v1Group)
		routers.PubSystem(v1Group)
	}

	{ // 需要认证
		v1Group.Use(authMiddleware.MiddlewareFunc(), middleware.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.RoleKey))
		routers.Base(v1Group, authMiddleware)
		routers.Dept(v1Group)
		routers.Dict(v1Group)
		routers.SysUser(v1Group)
		routers.Role(v1Group)
		routers.Config(v1Group)
		routers.UserCenter(v1Group)
		routers.Post(v1Group)
		routers.Menu(v1Group)
		routers.LoginLog(v1Group)
		routers.OperLog(v1Group)
	}
}

func OperLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" && deployed.EnabledDB {
			SetDBOperLog(c, c.ClientIP(), c.Writer.Status(), c.Request.RequestURI, c.Request.Method, time.Since(startTime))
		}
	}
}

// 写入操作日志表
// 该方法后续即将弃用
func SetDBOperLog(c *gin.Context, clientIP string, statusCode int, reqUri string, reqMethod string, latencyTime time.Duration) {
	menu := models.Menu{}
	menu.Path = reqUri
	menu.Action = reqMethod
	menuList, _ := menu.Get()
	sysOperLog := models.SysOperLog{
		OperIp:        clientIP,
		OperLocation:  deployed.IPLocation(clientIP),
		Status:        cast.ToString(statusCode),
		OperName:      jwtauth.UserName(c),
		RequestMethod: c.Request.Method,
		OperUrl:       reqUri,
	}

	if reqUri == "/login" {
		sysOperLog.BusinessType = "10"
		sysOperLog.Title = "用户登录"
		sysOperLog.OperName = "-"
	} else if strings.Contains(reqUri, "/api/v1/logout") {
		sysOperLog.BusinessType = "11"
	} else if strings.Contains(reqUri, "/api/v1/getCaptcha") {
		sysOperLog.BusinessType = "12"
		sysOperLog.Title = "验证码"
	} else {
		if reqMethod == "POST" {
			sysOperLog.BusinessType = "1"
		} else if reqMethod == "PUT" {
			sysOperLog.BusinessType = "2"
		} else if reqMethod == "DELETE" {
			sysOperLog.BusinessType = "3"
		}
	}
	sysOperLog.Method = reqMethod
	if len(menuList) > 0 {
		sysOperLog.Title = menuList[0].Title
	}
	b, _ := c.Get("body")
	sysOperLog.OperParam, _ = jsoniter.MarshalToString(b)
	sysOperLog.CreateBy = jwtauth.UserName(c)
	sysOperLog.OperTime = time.Now()
	sysOperLog.LatencyTime = (latencyTime).String()
	sysOperLog.UserAgent = c.Request.UserAgent()
	if c.Err() == nil {
		sysOperLog.Status = "0"
	} else {
		sysOperLog.Status = "1"
	}
	new(models.CallSysOperLog).Create(context.Background(), sysOperLog) // nolint: errcheck
}
