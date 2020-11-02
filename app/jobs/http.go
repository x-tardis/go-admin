package jobs

import (
	"log"
	"time"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/pkg"
)

var retryCount = 3

// HttpJob http job
type HttpJob struct {
	Base
}

// http 任务接口
func (h *HttpJob) Run() {
	startTime := time.Now()
	for count := 0; count < retryCount; count++ {
		str, err := pkg.Get(h.InvokeTarget)
		if err != nil {
			// 如果失败暂停一段时间重试
			log.Println("[ERROR] mission failed! ", err)
			log.Printf("[INFO] Retry after the task fails %d seconds! %s \n", time.Duration(count)*time.Second, str)
			time.Sleep(time.Duration(count) * time.Second)
			continue
		}
	}
	latencyTime := time.Since(startTime)
	// TODO: 待完善部分

	deployed.JobLogger.Infof("JobCore exec success %+v, spend: %v", h, latencyTime)
}
