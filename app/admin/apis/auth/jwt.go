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
	userIdKey    = "userId"
	usernameKey  = "username"
	roleIdKey    = "roleid"
	roleNameKey  = "rolename"
	roleKey      = "rolekey"
	dataScopeKey = "datascope"
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
					userIdKey:    v.UserId,
					usernameKey:  v.UserName,
					roleIdKey:    v.RoleId,
					roleNameKey:  v.RoleName,
					roleKey:      v.RoleKey,
					dataScopeKey: v.DataScope,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return infra.JWTIdentity{
				cast.ToInt(claims[userIdKey]),
				cast.ToString(claims[usernameKey]),
				cast.ToInt(claims[roleIdKey]),
				cast.ToString(claims[roleNameKey]),
				cast.ToString(claims[roleKey]),
				cast.ToString(claims[dataScopeKey]),
			}
		},
		Authenticator: authenticator,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			_, ok := data.(infra.JWTIdentity)
			return ok
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{"code": code, "msg": message})
		},
		LogoutResponse: logoutResponse,
		TokenLookup:    "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:  "Bearer",
		TimeFunc:       time.Now,
	})
}
