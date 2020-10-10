package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     deployed.SslConfig.Domain,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			return
		}
		c.Next()
	}
}
