package deployed

import "github.com/spf13/viper"

type Database struct {
	Driver string
	Source string
}

func ViperDatabase() *Database {
	cfg := viper.Sub("database")
	return &Database{
		Driver: cfg.GetString("driver"),
		Source: cfg.GetString("source"),
	}
}
