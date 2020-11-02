package jobs

import (
	"log"
)

// 需要将定义的struct 添加到字典中；
// 字典 key 可以配置到 自动任务 调用目标 中；
func initJob() {
	JobList = map[string]JobExec{
		"ExamplesOne": ExamplesOne{},
		// ...
	}
}

// 新添加的job 必须按照以下格式定义，并实现Exec函数
type ExamplesOne struct{}

func (ExamplesOne) Exec(arg interface{}) error {
	log.Printf("[INFO] JobCore ExamplesOne exec success, %v", arg)
	return nil
}
