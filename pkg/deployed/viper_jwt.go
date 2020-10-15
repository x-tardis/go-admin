package deployed

import (
	"time"

	"github.com/spf13/viper"

	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

func ViperJwtDefault() {
	viper.SetDefault("jwt.realm", "go-admin")
	viper.SetDefault("jwt.secretKey", "go-admin")
	viper.SetDefault("jwt.timeout", 24*7*time.Hour)
	viper.SetDefault("jwt.maxRefresh", 24*7*time.Hour)
}

func ViperJwt() *jwtauth.Config {
	cfg := viper.Sub("jwt")
	return &jwtauth.Config{
		cfg.GetString("realm"),
		cfg.GetString("secretKey"),
		cfg.GetDuration("timeout"),
		cfg.GetDuration("maxRefresh"),
	}
}
