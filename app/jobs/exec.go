package jobs

import (
	"time"

	"github.com/x-tardis/go-admin/deployed"
)

type ExecJob struct {
	Base
}

func (e *ExecJob) Run() {
	startTime := time.Now()
	var obj = jobList[e.InvokeTarget]
	if obj == nil {
		deployed.JobLogger.Warn(" ExecJob Run job nil", e)
		return
	}
	obj.Exec(e.Args)

	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	//TODO: 待完善部分
	//str := time.Now().Format(timeFormat) + " [INFO] JobCore " + string(e.EntryId) + "exec success , spend :" + latencyTime.String()
	//ws.SendAll(str)
	deployed.JobLogger.Info(time.Now().Format(timeFormat), " [INFO] JobCore ", e, "exec success , spend :", latencyTime)
}
