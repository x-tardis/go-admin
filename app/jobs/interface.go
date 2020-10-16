package jobs

type Job interface {
	Run()
	Expression() string
}

type JobExec interface {
	Exec(arg interface{}) error
}
