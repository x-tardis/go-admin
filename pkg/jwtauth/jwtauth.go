package jwtauth

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

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
	UserName  string
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
		return data.UserName
	}
	return ""
}

func RoleName(c *gin.Context) string {
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
