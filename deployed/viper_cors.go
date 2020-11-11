package deployed

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"
)

// ViperCorsDefault cors 默认值
func ViperCorsDefault() {
	viper.SetDefault("cors.allowOrigins", []string{})
	viper.SetDefault("cors.allowMethods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "UPDATE", "PATCH"})
	viper.SetDefault("cors.allowHeaders", []string{"Content-Type", "Content-Length", "AccessToken", "X-CSRF-Token", "Authorization", "Token", "X-Token"})
	viper.SetDefault("cors.allowCredentials", true)
	viper.SetDefault("cors.exposeHeaders", []string{"Content-Type", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"})
	viper.SetDefault("cors.maxAge", 300*time.Second)
}

// ViperCors cors config
func ViperCors() cors.Config {
	c := viper.Sub("cors")
	return cors.Config{
		AllowAllOrigins: c.GetBool("allowAllOrigins"),
		AllowOrigins:    c.GetStringSlice("allowOrigins"), // *
		//AllowOriginFunc: func(origin string) bool { return true },
		AllowMethods:           c.GetStringSlice("allowMethods"),  // "GET", "POST", "PUT", "DELETE", "OPTIONS", "UPDATE", "PATCH"
		AllowHeaders:           c.GetStringSlice("allowHeaders"),  // "Accept", "Authorization", "Content-Length", "X-CSRF-Token"
		AllowCredentials:       c.GetBool("allowCredentials"),     // true
		ExposeHeaders:          c.GetStringSlice("exposeHeaders"), // "Link"
		MaxAge:                 c.GetDuration("maxAge"),           // 300s
		AllowWildcard:          c.GetBool("allowWildcard"),
		AllowBrowserExtensions: c.GetBool("allowBrowserExtensions"),
		AllowWebSockets:        c.GetBool("allowWebSockets"),
		AllowFiles:             c.GetBool("allowFiles"),
	}
}
