package deployed

import "github.com/spf13/viper"

type Application struct {
	ReadTimeout   int
	WriterTimeout int
	Host          string
	Port          string
	Name          string
	JwtSecret     string
	Mode          string
	DemoMsg       string
	EnableDP      bool
}

var ApplicationConfig = new(Application)

func ViperApplication() *Application {
	cfg := viper.Sub("application")
	return &Application{
		ReadTimeout:   cfg.GetInt("readTimeout"),
		WriterTimeout: cfg.GetInt("writerTimeout"),
		Host:          cfg.GetString("host"),
		Port:          cfg.GetString("port"),
		Name:          cfg.GetString("name"),
		JwtSecret:     cfg.GetString("jwtSecret"),
		Mode:          cfg.GetString("mode"),
		DemoMsg:       cfg.GetString("demoMsg"),
		EnableDP:      cfg.GetBool("enabledp"),
	}
}
