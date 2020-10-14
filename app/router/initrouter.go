package router

import (
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"github.com/thinkgos/gin-middlewares/gzap"

	"github.com/x-tardis/go-admin/app/apis/auth"
	"github.com/x-tardis/go-admin/app/apis/ping"
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
	engine.Use(middleware.WithContextDb(middleware.GetGormFromConfig(deployed.Cfg)))
	engine.Use(
		gzap.Logger(deployed.RequestLogger.Desugar()), // logger
		gzap.Recovery(izap.Logger, false),             // recover
		OperLog(),                                     // 操作日志写入数据库
		middleware.NoCache(),                          // NoCache is a middleware function that appends headers
		middleware.Cors(),                             // 跨域处理
		middleware.Secure(),                           // Secure is a middleware function that appends security
	)
	// the jwt middleware
	authMiddleware, err := auth.NewJWTAuth(deployed.ApplicationConfig.JwtSecret)
	if err != nil {
		panic("JWT Init Error")
	}

	// 注册系统路由
	InitSysRouter(engine, authMiddleware)

	// 注册业务路由
	// TODO: 这里可存放业务路由，里边并无实际路由只有演示代码
	InitBusiness(engine, authMiddleware)

	return engine
}

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendAllService()

	g := r.Group("")

	g.GET("/", system.HelloWorld)
	g.GET("/info", ping.Ping)
	g.POST("/login", authMiddleware.LoginHandler)
	g.GET("/refresh_token", authMiddleware.RefreshHandler) // Refresh time can be longer than token timeout

	// 静态文件
	sysStaticFileRouter(g)
	// swagger
	routers.Swagger(g)

	// 需要认证
	wsGroup := r.Group("")
	wsGroup.Use(authMiddleware.MiddlewareFunc())
	{
		wsGroup.GET("/ws/:id/:channel", ws.WebsocketManager.WsClient)
		wsGroup.GET("/wslogout/:id/:channel", ws.WebsocketManager.UnWsClient)
	}

	v1 := r.Group("/api/v1")
	{ // 无需认证
		routers.NoCheckRoleBase(v1)

		routers.DB(v1)
		routers.SysTable(v1)
		routers.Public(v1)
		routers.SysSetting(v1)
	}

	{ // 需要认证
		v1.Use(authMiddleware.MiddlewareFunc(), middleware.NewAuthorizer(deployed.CasbinEnforcer, jwtauth.RoleKey))
		routers.Page(v1)
		routers.Base(v1, authMiddleware)
		routers.Dept(v1)
		routers.Dict(v1)
		routers.SysUser(v1)
		routers.Role(v1)
		routers.Config(v1)
		routers.UserCenter(v1)
		routers.Post(v1)
		routers.Menu(v1)
		routers.LoginLog(v1)
		routers.OperLog(v1)
	}
	return g
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
	sysOperLog := models.SysOperLog{}
	sysOperLog.OperIp = clientIP
	sysOperLog.OperLocation = deployed.IPLocation(clientIP)
	sysOperLog.Status = cast.ToString(statusCode)
	sysOperLog.OperName = jwtauth.UserName(c)
	sysOperLog.RequestMethod = c.Request.Method
	sysOperLog.OperUrl = reqUri
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
	_, _ = sysOperLog.Create()
}
