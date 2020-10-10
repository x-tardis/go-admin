package config

import "github.com/spf13/viper"

type Gen struct {
	DBName    string
	FrontPath string
}

var GenConfig = new(Gen)

func InitGen(cfg *viper.Viper) *Gen {
	return &Gen{
		DBName:    cfg.GetString("dbname"),
		FrontPath: cfg.GetString("frontpath"),
	}
}
