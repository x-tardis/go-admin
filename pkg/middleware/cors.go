package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    []string{"*"},
		// AllowOriginFunc: func(origin string) bool { return true },
		AllowMethods:           []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "UPDATE", "PATCH"},
		AllowHeaders:           []string{"Content-Type", "X-CSRF-Token", "Authorization", "authorization", "origin", "content-type", "accept"},
		AllowCredentials:       true,
		ExposeHeaders:          []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		MaxAge:                 300 * time.Second,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        true,
		AllowFiles:             false,
	})
}
