package dao

import (
	"github.com/spf13/viper"
	"github.com/x-tardis/go-admin/pkg/database"
	"github.com/x-tardis/go-admin/pkg/infra"
)

var DbConfig = database.Config{}

func ViperDatabase() database.Config {
	c := viper.Sub("database")
	cf := database.Config{
		Dialect:  c.GetString("dialect"),
		Username: c.GetString("username"),
		Password: c.GetString("password"),
		Protocol: c.GetString("protocol"),
		Host:     c.GetString("host"),
		Port:     c.GetString("port"),
		DbName:   c.GetString("dbName"),
		Extend:   make(map[string]string),
		LogMode:  c.GetBool("logMode"),
	}
	var extend []map[string]string

	err := c.UnmarshalKey("extend", &extend)
	infra.HandlerError(err)

	for _, v := range extend {
		key, ok1 := v["key"]
		value, ok2 := v["value"]
		if ok1 && ok2 && key != "" && value != "" {
			cf.Extend[key] = value
		}
	}
	return cf
}
