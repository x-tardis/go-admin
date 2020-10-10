package deployed

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"
)

// ViperCorsDefault cors 默认值
func ViperCorsDefault() {
	viper.SetDefault("cors.allowOrigins", []string{"*"})
	viper.SetDefault("cors.allowMethods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "UPDATE", "PATCH"})
	viper.SetDefault("cors.allowHeaders", []string{"Content-Type", "AccessToken", "X-CSRF-Token", "Authorization", "Token", "X-Token", "X-User-Id\""})
	viper.SetDefault("cors.allowCredentials", true)
	viper.SetDefault("cors.exposeHeaders", []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"})
	viper.SetDefault("cors.maxAge", 300*time.Second)
}

// ViperCors cors config
func ViperCors() cors.Config {
	return cors.Config{
		AllowOrigins: viper.GetStringSlice("cors.allowOrigins"), // *
		//AllowOriginFunc: func(origin string) bool { return true },
		AllowMethods:     viper.GetStringSlice("cors.allowMethods"),  // "GET", "POST", "PUT", "DELETE", "OPTIONS", "UPDATE", "PATCH"
		AllowHeaders:     viper.GetStringSlice("cors.allowHeaders"),  // "Accept", "Authorization", "Content-Length", "X-CSRF-Token"
		AllowCredentials: viper.GetBool("cors.allowCredentials"),     // true
		ExposeHeaders:    viper.GetStringSlice("cors.exposeHeaders"), // "Link"
		MaxAge:           viper.GetDuration("cors.maxAge"),           // 300s
	}
}
