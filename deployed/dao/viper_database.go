package dao

import (
	"github.com/spf13/viper"
	"github.com/x-tardis/go-admin/pkg/database"
)

var DbConfig = database.Config{}

func ViperDatabase() database.Config {
	c := viper.Sub("database")
	return database.Config{
		Dialect:  c.GetString("dialect"),
		Username: c.GetString("username"),
		Password: c.GetString("password"),
		Protocol: c.GetString("protocol"),
		Host:     c.GetString("host"),
		Port:     c.GetString("port"),
		DbName:   c.GetString("dbName"),
		Extend:   c.GetStringMapString("extend"),
		LogMode:  c.GetBool("logMode"),
	}
}
