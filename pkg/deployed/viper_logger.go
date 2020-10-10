package deployed

import "github.com/spf13/viper"

type Logger struct {
	Path       string
	Level      string
	Stdout     bool
	EnabledBUS bool
	EnabledREQ bool
	EnabledDB  bool
	EnabledJOB bool `default:"false"`
}

var LoggerConfig = new(Logger)

func ViperLogger() *Logger {
	cfg := viper.Sub("logger")
	return &Logger{
		Path:       cfg.GetString("path"),
		Level:      cfg.GetString("level"),
		Stdout:     cfg.GetBool("stdout"),
		EnabledBUS: cfg.GetBool("enabledbus"),
		EnabledREQ: cfg.GetBool("enabledreq"),
		EnabledDB:  cfg.GetBool("enableddb"),
		EnabledJOB: cfg.GetBool("enabledjob"),
	}
}
