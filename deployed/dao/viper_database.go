package dao

import (
	"github.com/spf13/viper"
	"github.com/x-tardis/go-admin/pkg/database"
)

var DbConfig = new(database.Database)

func ViperDatabase() *database.Database {
	c := viper.Sub("database")
	return &database.Database{
		Driver:    c.GetString("driver"),
		Source:    c.GetString("source"),
		EnableLog: c.GetBool("enableLog"),
	}
}
