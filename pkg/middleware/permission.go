package middleware

import (
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/tools"
)

// 权限检查中间件
func AuthCheckRole(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		v := jwt.ExtractClaims(c)
		// 检查权限
		res, err := enforcer.Enforce(v["rolekey"], c.Request.URL.Path, c.Request.Method)
		tools.HasError(err, "", 500)

		fmt.Printf("%s [INFO] %s %s %s \r\n",
			tools.GetCurrentTimeStr(),
			c.Request.Method,
			c.Request.URL.Path,
			v["rolekey"],
		)

		if res {
			c.Next()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "对不起，您没有该接口访问权限，请联系管理员",
			})
			c.Abort()
			return
		}
	}
}
