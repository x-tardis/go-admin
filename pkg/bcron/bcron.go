package bcron

import (
	"github.com/robfig/cron/v3"
)

var Cron = cron.New(cron.WithParser(cron.NewParser(cron.Second | cron.Minute |
	cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)))
