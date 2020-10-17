package jwtauth

import (
	"context"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type IdentityKey struct{}

// Config jwt 配置信息
type Config struct {
	Realm      string        `yaml:"realm" json:"realm"`
	SecretKey  string        `yaml:"secretKey" json:"secretKey"`
	Timeout    time.Duration `yaml:"timeout" json:"timeout"`
	MaxRefresh time.Duration `yaml:"maxRefresh" json:"maxRefresh"`
}

// Identities jwt identity
type Identities struct {
	UserId    int
	Username  string
	RoleId    int
	RoleName  string
	RoleKey   string
	DataScope string
}

func Identity(c *gin.Context) (Identities, bool) {
	data, exist := c.Get(jwt.IdentityKey)
	if exist {
		return data.(Identities), true
	}
	log.Println("[WARING] " + c.Request.Method + " " + c.Request.URL.Path + " 缺少 jwt.IdentityKey")
	return Identities{}, false
}

func UserId(c *gin.Context) int {
	if data, ok := Identity(c); ok {
		return data.UserId
	}
	return 0
}

func UserIdStr(c *gin.Context) string {
	return cast.ToString(UserId(c))
}

func UserName(c *gin.Context) string {
	if data, ok := Identity(c); ok {
		return data.Username
	}
	return ""
}

func RoleName(c *gin.Context) string {
	if data, ok := Identity(c); ok {
		return data.RoleName
	}
	return ""
}

func RoleKey(c *gin.Context) string {
	if data, ok := Identity(c); ok {
		return data.RoleKey
	}
	return ""
}

func RoleId(c *gin.Context) int {
	if data, ok := Identity(c); ok {
		return data.RoleId
	}
	return 0
}

func DataScope(c *gin.Context) string {
	if data, ok := Identity(c); ok {
		return data.DataScope
	}
	return ""
}

func FromIdentity(ctx context.Context) Identities {
	return ctx.Value(IdentityKey{}).(Identities)
}

func FromUserId(ctx context.Context) int {
	return FromIdentity(ctx).UserId
}
func FromUserIdStr(ctx context.Context) string {
	return cast.ToString(FromIdentity(ctx).UserId)
}

func FromUserName(ctx context.Context) string {
	return FromIdentity(ctx).Username
}

func FromRoleName(ctx context.Context) string {
	return FromIdentity(ctx).RoleName
}

func FromRoleKey(ctx context.Context) string {
	return FromIdentity(ctx).RoleKey
}

func FromRoleId(ctx context.Context) int {
	return FromIdentity(ctx).RoleId
}

func FromDataScope(ctx context.Context) string {
	return FromIdentity(ctx).DataScope
}
