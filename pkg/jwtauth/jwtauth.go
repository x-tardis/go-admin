package jwtauth

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"
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

func CasbinSubject(c *gin.Context) string {
	return FromRoleKey(gcontext.Context(c))
}
