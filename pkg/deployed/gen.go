package deployed

import "github.com/spf13/viper"

type Gen struct {
	DBName    string
	FrontPath string
}

var GenConfig = new(Gen)

func ViperGen() *Gen {
	cfg := viper.Sub("gen")
	return &Gen{
		DBName:    cfg.GetString("dbname"),
		FrontPath: cfg.GetString("frontpath"),
	}
}
