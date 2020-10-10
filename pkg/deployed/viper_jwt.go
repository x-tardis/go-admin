package deployed

import (
	"github.com/spf13/viper"
)

type Jwt struct {
	Secret  string
	Timeout int64
}

var JwtConfig = new(Jwt)

func ViperJwt() *Jwt {
	cfg := viper.Sub("jwt")
	return &Jwt{
		Secret:  cfg.GetString("secret"),
		Timeout: cfg.GetInt64("timeout"),
	}
}
