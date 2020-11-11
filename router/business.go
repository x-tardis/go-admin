package router

import (
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/routers"
)

var businessMu = sync.RWMutex{}

// 业务公共开放接口
var businessPublic []func(r gin.IRouter)

// 业务有授权,有RBAC角色控制
var businessAuthRbac []func(r gin.IRouter)

func init() {
	RegisterBusinessPublic(
		routers.PubSysFileInfo,
		routers.PubSysFileDir,
	)
	RegisterBusinessAuthRbac(
		routers.SysJobRouter,
		routers.SysContent,
		routers.SysCategory,
	)
}

func RegisterBusinessPublic(f ...func(r gin.IRouter)) {
	businessMu.Lock()
	businessPublic = append(businessPublic, f...)
	businessMu.Unlock()
}

func RegisterBusinessAuthRbac(f ...func(r gin.IRouter)) {
	businessMu.Lock()
	businessAuthRbac = append(businessAuthRbac, f...)
	businessMu.Unlock()
}
