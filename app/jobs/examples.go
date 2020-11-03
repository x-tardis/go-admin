package jobs

import (
	"log"
)

// 新添加的job 必须按照以下格式定义，并实现Exec函数
type ExamplesOne struct{}

func (ExamplesOne) Exec(arg interface{}) error {
	log.Printf("[INFO] JobCore ExamplesOne exec success, %v", arg)
	return nil
}
