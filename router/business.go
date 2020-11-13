package router

import (
	"sync"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/routers"
)

var businessMu = sync.RWMutex{}

var custom = func(*gin.Engine, *jwt.GinJWTMiddleware) {}

// 业务公共开放接口
var businessPublic []func(v1 gin.IRouter)

// 业务有授权,有RBAC角色控制
var businessAuthRbac []func(v1 gin.IRouter)

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

func RegisterCustom(f func(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware)) {
	if f == nil {
		panic("router: custom function is nil")
	}
	custom = f
}

func RegisterBusinessPublic(f ...func(v1 gin.IRouter)) {
	businessMu.Lock()
	businessPublic = append(businessPublic, f...)
	businessMu.Unlock()
}

func RegisterBusinessAuthRbac(f ...func(v1 gin.IRouter)) {
	businessMu.Lock()
	businessAuthRbac = append(businessAuthRbac, f...)
	businessMu.Unlock()
}
