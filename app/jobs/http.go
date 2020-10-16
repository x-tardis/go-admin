package jobs

import (
	"fmt"
	"time"

	"github.com/x-tardis/go-admin/pkg"
	"github.com/x-tardis/go-admin/pkg/deployed"
)

var retryCount = 3

// 任务类型 http
type HttpJob struct {
	Base
}

// http 任务接口
func (h *HttpJob) Run() {
	startTime := time.Now()
	var count = 0
	/* 循环 */
LOOP:
	if count < retryCount {
		/* 跳过迭代 */
		str, err := pkg.Get(h.InvokeTarget)
		if err != nil {
			// 如果失败暂停一段时间重试
			fmt.Println(time.Now().Format(timeFormat), " [ERROR] mission failed! ", err)
			fmt.Printf(time.Now().Format(timeFormat)+" [INFO] Retry after the task fails %d seconds! %s \n", time.Duration(count)*time.Second, str)
			time.Sleep(time.Duration(count) * time.Second)
			goto LOOP
		}
		count++
	}
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	// TODO: 待完善部分

	deployed.JobLogger.Info(time.Now().Format(timeFormat), " [INFO] JobCore ", h, "exec success , spend :", latencyTime)
}
