package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func NewAuthorizer(e *casbin.SyncedEnforcer, subject func(c *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, err := e.Enforce(subject(c), c.Request.URL.Path, c.Request.Method)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Permission validation errors occur!",
			})
		} else if !allowed {
			// the 403 Forbidden to the client
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": http.StatusForbidden,
				"msg":  "对不起，您没有该接口访问权限，请联系管理员", // "Permission denied!",
			})
			return
		}
		c.Next()
	}
}
