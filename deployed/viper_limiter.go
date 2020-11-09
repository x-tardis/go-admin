package deployed

import (
	"github.com/spf13/viper"
)

func ViperLimiterDefault() {
	viper.SetDefault("rate.limit", float64(10000))
}

func ViperLimiter() float64 {
	return viper.GetFloat64("rate.limit")
}
