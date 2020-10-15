package deployed

import "github.com/spf13/viper"

type Ssl struct {
	KeyStr string
	Pem    string
	Enable bool
	Domain string
}

func ViperSsl() *Ssl {
	c := viper.Sub("ssl")
	return &Ssl{
		KeyStr: c.GetString("key"),
		Pem:    c.GetString("pem"),
		Enable: c.GetBool("enable"),
		Domain: c.GetString("domain"),
	}
}
