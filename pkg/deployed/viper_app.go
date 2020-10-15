package deployed

import (
	"github.com/spf13/viper"
)

type Application struct {
	Mode          string
	Name          string
	Host          string
	Port          string
	ReadTimeout   int
	WriterTimeout int
	EnableDP      bool
}

var AppConfig = new(Application)

func ViperApplication() *Application {
	return &Application{
		viper.GetString("mode"),
		viper.GetString("name"),
		viper.GetString("host"),
		viper.GetString("port"),
		viper.GetInt("readTimeout"),
		viper.GetInt("writerTimeout"),
		viper.GetBool("enabledp"),
	}
}
