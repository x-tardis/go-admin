package router

import (
	"sync"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/routers"
)

// business 业务接口
type business struct {
	rw       sync.RWMutex
	custom   []func(*gin.Engine, *jwt.GinJWTMiddleware)
	public   []func(v1 gin.IRouter)
	authRBAC []func(v1 gin.IRouter)
}

var businesses = &business{
	public: []func(v1 gin.IRouter){
		routers.PubSysFileInfo,
		routers.PubSysFileDir,
	},
	authRBAC: []func(v1 gin.IRouter){
		routers.SysJobRouter,
		routers.SysContent,
		routers.SysCategory,
	},
}

func RegisterBusinessCustom(f ...func(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware)) {
	businesses.rw.Lock()
	businesses.custom = append(businesses.custom, f...)
	businesses.rw.Unlock()
}

// RegisterBusinessPublic 注册公开业务接口,属 "/api/v1" 组
func RegisterBusinessPublic(f ...func(v1 gin.IRouter)) {
	businesses.rw.Lock()
	businesses.public = append(businesses.public, f...)
	businesses.rw.Unlock()
}

// RegisterBusinessPublic 注册业务接口,由RBAC授权控制,属 "/api/v1" 组
func RegisterBusinessAuthRbac(f ...func(v1 gin.IRouter)) {
	businesses.rw.Lock()
	businesses.authRBAC = append(businesses.authRBAC, f...)
	businesses.rw.Unlock()
}
