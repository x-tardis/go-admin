package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/sysfiledir"
	"github.com/x-tardis/go-admin/app/apis/sysfileinfo"
)

// 无需认证的路由代码
func SysFileInfo(v1 gin.IRouter) {
	v1.GET("/sysfileinfoList", sysfileinfo.GetSysFileInfoList)
	r := v1.Group("/sysfileinfo")
	{
		r.GET("/:id", sysfileinfo.GetSysFileInfo)
		r.POST("", sysfileinfo.InsertSysFileInfo)
		r.PUT("", sysfileinfo.UpdateSysFileInfo)
		r.DELETE("/:id", sysfileinfo.DeleteSysFileInfo)
	}
}

// 无需认证的路由代码
func SysFileDir(v1 gin.IRouter) {
	v1.GET("/sysfiledirList", sysfiledir.GetSysFileDirList)
	r := v1.Group("/sysfiledir")
	{
		r.GET("/:id", sysfiledir.GetSysFileDir)
		r.POST("", sysfiledir.InsertSysFileDir)
		r.PUT("", sysfiledir.UpdateSysFileDir)
		r.DELETE("/:id", sysfiledir.DeleteSysFileDir)
	}
}
