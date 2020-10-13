package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// 权限检查中间件
func AuthCheckRole(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := jwtauth.RoleKey(c)
		// 检查权限
		res, err := enforcer.Enforce(role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			servers.Fail(c, 500, err.Error())
			return
		}

		fmt.Printf("%s [INFO] %s %s %s \r\n",
			time.Now().Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			role,
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
