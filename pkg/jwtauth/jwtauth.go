package jwtauth

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"
)

// IdentityKey identity key
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
	DeptId    int
	PostId    int
	RoleId    int
	RoleName  string
	RoleKey   string
	DataScope string
}

// FromIdentity get identity from context
func FromIdentity(ctx context.Context) Identities {
	return ctx.Value(IdentityKey{}).(Identities)
}

// FromUserId get userId from context
func FromUserId(ctx context.Context) int {
	return FromIdentity(ctx).UserId
}

// FromUserIdStr get userId string from context
func FromUserIdStr(ctx context.Context) string {
	return cast.ToString(FromIdentity(ctx).UserId)
}

// FromUserName get username from context
func FromUserName(ctx context.Context) string {
	return FromIdentity(ctx).Username
}

// FromDeptId get deptId from context
func FromDeptId(ctx context.Context) int {
	return FromIdentity(ctx).DeptId
}

// FromPostId get postId from context
func FromPostId(ctx context.Context) int {
	return FromIdentity(ctx).PostId
}

// FromRoleName get roleName from context
func FromRoleName(ctx context.Context) string {
	return FromIdentity(ctx).RoleName
}

// FromRoleKey get roleKey from context
func FromRoleKey(ctx context.Context) string {
	return FromIdentity(ctx).RoleKey
}

// FromRoleId get roleId from context
func FromRoleId(ctx context.Context) int {
	return FromIdentity(ctx).RoleId
}

// FromDataScope get dataScope from context
func FromDataScope(ctx context.Context) string {
	return FromIdentity(ctx).DataScope
}

// CasbinSubject get casbin subject from context, which is roleKey.
func CasbinSubject(c *gin.Context) string {
	return FromRoleKey(gcontext.Context(c))
}
