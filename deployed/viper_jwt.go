package deployed

import (
	"time"

	"github.com/spf13/viper"

	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// TODO: bug 当jwt token过期时,页面并未有跳转到登录页面
func ViperJwtDefault() {
	viper.SetDefault("jwt.realm", "go-admin")
	viper.SetDefault("jwt.secretKey", "go-admin")
	viper.SetDefault("jwt.timeout", 24*7*time.Hour)
	viper.SetDefault("jwt.maxRefresh", 24*7*time.Hour)
}

func ViperJwt() jwtauth.Config {
	cfg := viper.Sub("jwt")
	return jwtauth.Config{
		Realm:      cfg.GetString("realm"),
		SecretKey:  cfg.GetString("secretKey"),
		Timeout:    cfg.GetDuration("timeout"),
		MaxRefresh: cfg.GetDuration("maxRefresh"),
	}
}
