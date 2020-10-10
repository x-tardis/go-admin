package deployed

import (
	"github.com/robfig/cron/v3"
)

// NewCron returns a Cron with the seconds field enabled.
func NewCron() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
