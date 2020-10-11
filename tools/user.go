package tools

import (
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetUserId(c *gin.Context) int {
	data := jwt.ExtractClaims(c)
	if data["userId"] != nil {
		return int((data["userId"]).(float64))
	}
	fmt.Println(GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserId 缺少identity")
	return 0
}

func GetUserIdStr(c *gin.Context) string {
	data := jwt.ExtractClaims(c)
	if data["userId"] != nil {
		return cast.ToString(int64((data["userId"]).(float64)))
	}
	fmt.Println(GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserIdStr 缺少identity")
	return ""
}

func GetUserName(c *gin.Context) string {
	data := jwt.ExtractClaims(c)
	if data["UsernameKey"] != nil {
		return (data["UsernameKey"]).(string)
	}
	fmt.Println(GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserName 缺少nice")
	return ""
}

func GetRoleName(c *gin.Context) string {
	data := jwt.ExtractClaims(c)
	if data["rolekey"] != nil {
		return (data["rolekey"]).(string)
	}
	fmt.Println(GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleName 缺少rolekey")
	return ""
}

func GetRoleId(c *gin.Context) int {
	data := jwt.ExtractClaims(c)
	if data["roleid"] != nil {
		i := int((data["roleid"]).(float64))
		return i
	}
	fmt.Println(GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleId 缺少roleid")
	return 0
}
