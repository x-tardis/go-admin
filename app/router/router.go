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
	"github.com/thinkgos/gin-middlewares/authj"
	"github.com/thinkgos/gin-middlewares/expvar"
	"github.com/thinkgos/gin-middlewares/gzap"
	"github.com/thinkgos/gin-middlewares/requestid"
	"github.com/thinkgos/sharp/gin/gcontext"

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
		requestid.RequestID(),
		gzap.Logger(deployed.RequestLogger.Desugar()), // logger
		gzap.Recovery(izap.Logger, false),             // recover
		OperLog(),                                     // 操作日志写入数据库
		middleware.NoCache(),                          // NoCache is a middleware function that appends headers
		cors.New(*deployed.CorsConfig),                // 跨域处理
		middleware.Secure(),                           // Secure is a middleware function that appends security
	)
	// the jwt middleware
	authMiddleware, err := system.NewJWTAuth(deployed.JwtConfig)
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
	engine.GET("/debug/vars", expvar.Handler())
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
		v1Group.Use(authMiddleware.MiddlewareFunc(), authj.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.CasbinSubject))
		routers.Base(v1Group, authMiddleware)
		routers.Dept(v1Group)
		routers.Dict(v1Group)
		routers.User(v1Group)
		routers.Role(v1Group)
		routers.Config(v1Group)
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
			SetDBOperLog(c, c.Writer.Status(), time.Since(startTime))
		}
	}
}

// 写入操作日志表
// 该方法后续即将弃用
func SetDBOperLog(c *gin.Context, statusCode int, latencyTime time.Duration) {
	reqUri := c.Request.RequestURI
	reqMethod := c.Request.Method
	menuList, _ := models.CMenu.Query(gcontext.Context(c), models.MenuQueryParam{
		Path:   reqUri,
		Action: reqMethod,
	})
	clientIP := c.ClientIP()
	username := jwtauth.FromUserName(gcontext.Context(c))
	sysOperLog := models.OperLog{
		OperIp:        clientIP,
		OperLocation:  deployed.IPLocation(clientIP),
		Status:        cast.ToString(statusCode),
		OperName:      username,
		RequestMethod: reqMethod,
		OperUrl:       reqUri,
		Creator:       username,
		OperTime:      time.Now(),
		LatencyTime:   latencyTime.String(),
		UserAgent:     c.Request.UserAgent(),
	}

	if reqUri == "/login" {
		sysOperLog.BusinessType = "10"
		sysOperLog.Title = "用户登录"
		sysOperLog.OperName = "-"
	} else if strings.Contains(reqUri, "/api/v1/logout") {
		sysOperLog.BusinessType = "11"
	} else if strings.Contains(reqUri, "/api/v1/captcha") {
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
	if c.Err() == nil {
		sysOperLog.Status = "0"
	} else {
		sysOperLog.Status = "1"
	}
	models.COperLog.Create(context.Background(), sysOperLog) // nolint: errcheck
}
