package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/sysfiledir"
	"github.com/x-tardis/go-admin/app/apis/sysfileinfo"
)

// 无需认证的路由代码
func PubSysFileInfo(v1 gin.IRouter) {
	r := v1.Group("/sysfileinfo")
	{
		r.GET("", sysfileinfo.GetSysFileInfoList)
		r.GET("/:id", sysfileinfo.GetSysFileInfo)
		r.POST("", sysfileinfo.InsertSysFileInfo)
		r.PUT("", sysfileinfo.UpdateSysFileInfo)
		r.DELETE("/:ids", sysfileinfo.DeleteSysFileInfo)
	}
}

// 无需认证的路由代码
func PubSysFileDir(v1 gin.IRouter) {
	r := v1.Group("/sysfiledir")
	{
		r.GET("", sysfiledir.GetSysFileDirList)
		r.GET("/:id", sysfiledir.GetSysFileDir)
		r.POST("", sysfiledir.InsertSysFileDir)
		r.PUT("", sysfiledir.UpdateSysFileDir)
		r.DELETE("/:ids", sysfiledir.DeleteSysFileDir)
	}
}
