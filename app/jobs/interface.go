package jobs

type Job interface {
	Expression() string
	Run()
}

type JobExec interface {
	Exec(arg interface{}) error
}
