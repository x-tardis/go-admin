package deployed

import (
	"math"

	"github.com/spf13/viper"
)

func ViperLimiterDefault() {
	viper.SetDefault("rate.limit", math.MaxFloat64)
}

func ViperLimiter() float64 {
	return viper.GetFloat64("rate.limit")
}
