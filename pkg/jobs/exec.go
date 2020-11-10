package jobs

import (
	"time"

	"github.com/x-tardis/go-admin/deployed"
)

// ExecJob exec job
type ExecJob struct {
	Base
}

func (e *ExecJob) Run() {
	obj := jobList[e.InvokeTarget]
	if obj == nil {
		deployed.JobLogger.Warn(" exec job is nil", e)
		return
	}

	startTime := time.Now()
	obj.Exec(e.Args)
	latencyTime := time.Since(startTime)

	//TODO: 待完善部分
	//ws.SendBroadcast(str)

	deployed.JobLogger.Infof("JobCore exec success %+v, spend: %v", e, latencyTime)
}
