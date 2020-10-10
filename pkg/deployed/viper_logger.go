package deployed

import "github.com/spf13/viper"

type Log struct {
	Path       string
	Level      string
	Stdout     bool
	EnabledBUS bool
	EnabledREQ bool
	EnabledDB  bool
	EnabledJOB bool `default:"false"`
}

var LoggerConfig = new(Log)

func ViperLogger() *Log {
	cfg := viper.Sub("logger")
	return &Log{
		Path:       cfg.GetString("path"),
		Level:      cfg.GetString("level"),
		Stdout:     cfg.GetBool("stdout"),
		EnabledBUS: cfg.GetBool("enabledbus"),
		EnabledREQ: cfg.GetBool("enabledreq"),
		EnabledDB:  cfg.GetBool("enableddb"),
		EnabledJOB: cfg.GetBool("enabledjob"),
	}
}
