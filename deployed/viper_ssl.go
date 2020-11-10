package deployed

import "github.com/spf13/viper"

type Ssl struct {
	KeyStr string
	Pem    string
	Enable bool
	Domain string
}

func ViperSsl() *Ssl {
	if c := viper.Sub("ssl"); c != nil {
		return &Ssl{
			c.GetString("key"),
			c.GetString("pem"),
			c.GetBool("enable"),
			c.GetString("domain"),
		}
	}
	return &Ssl{}
}
