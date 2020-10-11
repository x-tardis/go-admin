package auth

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/pkg/infra"
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
			if v, ok := data.(infra.JWTIdentity); ok {
				return jwt.MapClaims{
					UserIdKey:    v.UserId,
					UsernameKey:  v.UserName,
					RoleIdKey:    v.RoleId,
					RoleNameKey:  v.RoleName,
					RoleKey:      v.RoleKey,
					DataScopeKey: v.DataScope,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			// return map[string]interface{}{
			// 	"UserId":      claims[UserIdKey],
			// 	"IdentityKey": claims[UserIdKey],
			// 	"UserName":    claims[UsernameKey],
			// 	"RoleIds":     claims[RoleIdKey],
			// 	"RoleName":    claims[RoleNameKey],
			// 	"RoleKey":     claims[RoleKey],
			// 	"DataScope":   claims[DataScopeKey],
			// }
			return infra.JWTIdentity{
				cast.ToInt(claims[UserIdKey]),
				cast.ToString(claims[UsernameKey]),
				cast.ToInt(claims[RoleIdKey]),
				cast.ToString(claims[RoleNameKey]),
				cast.ToString(claims[RoleKey]),
				cast.ToString(claims[DataScopeKey]),
			}
		},
		Authenticator: authenticator,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(map[string]interface{}); ok {
			// 	u, _ := v["user"].(models.SysUser)
			// 	r, _ := v["role"].(models.SysRole)
			// 	c.Set("role", r.RoleName)
			// 	c.Set("roleIds", r.RoleId)
			// 	c.Set("userId", u.UserId)
			// 	c.Set("userName", u.UserName)
			// 	c.Set("dataScope", r.DataScope)
			// 	return true
			// }
			// return false
			_, ok := data.(infra.JWTIdentity)
			return ok
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{"code": code, "msg": message})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
