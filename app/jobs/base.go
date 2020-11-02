package jobs

type Job interface {
	Expression() string
	Run()
	SetEntryId(id int)
}

type JobExec interface {
	Exec(arg interface{}) error
}

// Base base job
type Base struct {
	JobId          uint   // 数据库id
	EntryId        int    // cron id
	Name           string // 名称
	InvokeTarget   string // 调用目标名
	CronExpression string // cron linux表达式
	Args           string // 回调参数
}

// Run implement Job interface
func (*Base) Run() {}

// Expression implement Job interface
func (b *Base) Expression() string { return b.CronExpression }

// SetEntryId implement Job interface
func (b *Base) SetEntryId(id int) { b.EntryId = id }
