package tools

import (
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/pkg/infra"
)

func GetJWTIdentity(c *gin.Context) (infra.JWTIdentity, bool) {
	data, exist := c.Get(jwt.IdentityKey)
	if exist {
		return data.(infra.JWTIdentity), true
	}
	return infra.JWTIdentity{}, false
}

func GetUserId(c *gin.Context) int {
	data, ok := GetJWTIdentity(c)
	if ok {
		return data.UserId
	}
	fmt.Println(GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserId 缺少identity")
	return 0
}

func GetUserIdStr(c *gin.Context) string {
	return cast.ToString(GetUserId(c))
}

func GetUserName(c *gin.Context) string {
	data := jwt.ExtractClaims(c)
	if data["username"] != nil {
		return (data["username"]).(string)
	}
	fmt.Println(GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserName 缺少username")
	return ""
}

func GetRoleName(c *gin.Context) string {
	data, ok := GetJWTIdentity(c)
	if ok {
		return data.RoleKey
	}
	fmt.Println(GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleName 缺少rolekey")
	return ""
}

func GetRoleId(c *gin.Context) int {
	data, ok := GetJWTIdentity(c)
	if ok {
		return data.RoleId
	}
	fmt.Println(GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleId 缺少roleid")
	return 0
}
