package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	xrequestid "github.com/thinkgos/http-middlewares/requestid"
)

func WriteResponseHeaderID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(xrequestid.RequestIDHeader, requestid.FromRequestID(c))
		c.Next()
	}
}
