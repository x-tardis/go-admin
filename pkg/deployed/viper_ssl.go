package deployed

import "github.com/spf13/viper"

type Ssl struct {
	KeyStr string
	Pem    string
	Enable bool
	Domain string
}

func ViperSsl() *Ssl {
	cfg := viper.Sub("ssl")
	if cfg == nil {
		cfg = viper.New()
	}
	return &Ssl{
		KeyStr: cfg.GetString("key"),
		Pem:    cfg.GetString("pem"),
		Enable: cfg.GetBool("enable"),
		Domain: cfg.GetString("domain"),
	}
}
