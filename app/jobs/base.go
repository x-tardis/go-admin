package jobs

type Base struct {
	InvokeTarget   string
	Name           string
	JobId          uint
	EntryId        int
	CronExpression string
	Args           string
}

func (*Base) Run() {}

func (b *Base) Expression() string { return b.CronExpression }
