package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
)

func WriteResponseHeaderID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Request-ID", requestid.FromRequestID(c))
		c.Next()
	}
}
