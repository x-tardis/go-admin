package router

import (
	"mime"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/monitor"
	"github.com/x-tardis/go-admin/app/apis/ping"
	"github.com/x-tardis/go-admin/app/apis/public"
	"github.com/x-tardis/go-admin/app/apis/system"
	"github.com/x-tardis/go-admin/app/apis/system/dict"
	. "github.com/x-tardis/go-admin/app/apis/tools"
	_ "github.com/x-tardis/go-admin/docs"
	"github.com/x-tardis/go-admin/pkg/ws"
)

func sysBaseRouter(r *gin.RouterGroup) {
	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendAllService()

	r.GET("/", system.HelloWorld)
	r.GET("/info", ping.Ping)
}

func sysStaticFileRouter(r *gin.RouterGroup) {
	mime.AddExtensionType(".js", "application/javascript")
	r.Static("/static", "./static")
	r.Static("/form-generator", "./static/form-generator")
}

func sysNoCheckRoleRouter(r gin.IRouter) {
	v1 := r.Group("/api/v1")

	v1.GET("/monitor/server", monitor.ServerInfo)
	v1.GET("/getCaptcha", system.GenerateCaptchaHandler)
	v1.GET("/gen/preview/:tableId", Preview)
	v1.GET("/gen/toproject/:tableId", GenCodeV3)
	v1.GET("/gen/todb/:tableId", GenMenuAndApi)
	v1.GET("/gen/tabletree", GetSysTablesTree)
	v1.GET("/menuTreeselect", system.GetMenuTreeelect)
	v1.GET("/dict/databytype/:dictType", dict.GetDictDataByDictType)

	registerDBRouter(v1)
	registerSysTableRouter(v1)
	registerPublicRouter(v1)
	registerSysSettingRouter(v1)
}

func registerDBRouter(api *gin.RouterGroup) {
	db := api.Group("/db")
	{
		db.GET("/tables/page", GetDBTableList)
		db.GET("/columns/page", GetDBColumnList)
	}
}

func registerSysTableRouter(v1 gin.IRouter) {
	systables := v1.Group("/sys/tables")
	{
		systables.GET("/page", GetSysTableList)
		tablesinfo := systables.Group("/info")
		{
			tablesinfo.POST("", InsertSysTable)
			tablesinfo.PUT("", UpdateSysTable)
			tablesinfo.DELETE("/:tableId", DeleteSysTables)
			tablesinfo.GET("/:tableId", GetSysTables)
			tablesinfo.GET("", GetSysTablesInfo)
		}
	}
}

func registerSysSettingRouter(v1 gin.IRouter) {
	setting := v1.Group("/setting")
	{
		setting.GET("", system.GetSetting)
		setting.POST("", system.CreateSetting)
		setting.GET("/serverInfo", monitor.ServerInfo)
	}
}

func registerPublicRouter(v1 gin.IRouter) {
	p := v1.Group("/public")
	{
		p.POST("/uploadFile", public.UploadFile)
	}
}
