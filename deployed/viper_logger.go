package deployed

import (
	"github.com/spf13/viper"

	"github.com/x-tardis/go-admin/pkg/izap"
)

func ViperLogger() izap.Config {
	c := viper.Sub("logger")
	return izap.Config{
		Level:       c.GetString("level"),
		Format:      c.GetString("format"),
		EncodeLevel: c.GetString("encodeLevel"),
		Adapter:     c.GetString("adapter"),
		Stack:       c.GetBool("stack"),
		Path:        c.GetString("path"),

		FileName:   c.GetString("fileName"),
		MaxSize:    c.GetInt("maxSize"),
		MaxAge:     c.GetInt("maxAge"),
		MaxBackups: c.GetInt("maxBackups"),
		LocalTime:  c.GetBool("localTime"),
		Compress:   c.GetBool("compress"),
	}
}
