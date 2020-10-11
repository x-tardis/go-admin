package auth

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/admin/models"
)

const (
	UserIdKey    = "userId"
	UsernameKey  = "username"
	RoleIdKey    = "roleid"
	RoleNameKey  = "rolename"
	RoleKey      = "rolekey"
	DataScopeKey = "datascope"
)

func NewJWTAuth(key string) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte(key),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(map[string]interface{}); ok {
				u, _ := v["user"].(models.SysUser)
				r, _ := v["role"].(models.SysRole)
				return jwt.MapClaims{
					UserIdKey:    u.UserId,
					UsernameKey:  u.Username,
					RoleIdKey:    r.RoleId,
					RoleNameKey:  r.RoleName,
					RoleKey:      r.RoleKey,
					DataScopeKey: r.DataScope,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return map[string]interface{}{
				"UserId":      claims[UserIdKey],
				"IdentityKey": claims[UserIdKey],
				"UserName":    claims[UsernameKey],
				"RoleIds":     claims[RoleIdKey],
				"RoleName":    claims[RoleNameKey],
				"RoleKey":     claims[RoleKey],
				"DataScope":   claims[DataScopeKey],
			}
		},
		Authenticator: authenticator,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(map[string]interface{}); ok {
				u, _ := v["user"].(models.SysUser)
				r, _ := v["role"].(models.SysRole)
				c.Set("role", r.RoleName)
				c.Set("roleIds", r.RoleId)
				c.Set("userId", u.UserId)
				c.Set("userName", u.UserName)
				c.Set("dataScope", r.DataScope)
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{"code": code, "msg": message})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
