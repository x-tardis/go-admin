package jobs

type Job interface {
	Expression() string
	Run()
	SetEntryId(id int)
}

type JobExec interface {
	Exec(arg interface{}) error
}
