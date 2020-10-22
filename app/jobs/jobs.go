package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/models"
)

var timeFormat = "2006-01-02 15:04:05"

var jobList map[string]JobExec

// 初始化
func Setup() {
	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Starting...")

	deployed.Cron = deployed.NewCron()

	jobList, err := models.CJob.Query(context.Background())
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore init error", err)
	}
	if len(jobList) == 0 {
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore total:0")
	}

	err = models.CJob.RemoveAllEntryID(context.Background())
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore remove entry_id error", err)
	}

	for i := 0; i < len(jobList); i++ {
		sysJob := models.Job{}
		if jobList[i].JobType == 1 {
			j := &HttpJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].JobName

			sysJob.EntryId, err = AddJob(j)
		} else if jobList[i].JobType == 2 {
			j := &ExecJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].JobName
			j.Args = jobList[i].Args
			sysJob.EntryId, err = AddJob(j)
		}
		err = models.CJob.Update(context.Background(), jobList[i].JobId, sysJob)
	}

	// 其中任务
	deployed.Cron.Start()
	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore start success.")
	// 关闭任务
	defer deployed.Cron.Stop()
	select {}
}

// 添加任务 AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(job Job) (int, error) {
	if job == nil {
		fmt.Println("unknown")
		return 0, nil
	}
	id, err := deployed.Cron.AddJob(job.Expression(), job)
	return int(id), err
}

// 移除任务
func Remove(entryID int) chan bool {
	ch := make(chan bool)
	go func() {
		deployed.Cron.Remove(cron.EntryID(entryID))
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Remove success ,info entryID :", entryID)
		ch <- true
	}()
	return ch
}

// 任务停止
func Stop() chan bool {
	ch := make(chan bool)
	go func() {
		deployed.Cron.Stop()
		ch <- true
	}()
	return ch
}
