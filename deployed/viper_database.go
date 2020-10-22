package deployed

import (
	"github.com/spf13/viper"
	"github.com/x-tardis/go-admin/pkg/database"
)

func ViperDatabase() *database.Database {
	cfg := viper.Sub("database")
	return &database.Database{
		Driver: cfg.GetString("driver"),
		Source: cfg.GetString("source"),
	}
}
