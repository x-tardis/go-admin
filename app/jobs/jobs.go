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
func Startup() {
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

	for _, v := range jobList {
		var job Job
		if v.JobType == 1 {
			job = &HttpJob{
				Base{
					InvokeTarget:   v.InvokeTarget,
					Name:           v.JobName,
					JobId:          v.JobId,
					EntryId:        0,
					CronExpression: v.CronExpression,
					Args:           "",
				},
			}
		} else if v.JobType == 2 {
			job = &ExecJob{
				Base{
					InvokeTarget:   v.InvokeTarget,
					Name:           v.JobName,
					JobId:          v.JobId,
					EntryId:        0,
					CronExpression: v.CronExpression,
					Args:           v.Args,
				},
			}
		}
		entryId, _ := AddJob(job)
		err = models.CJob.UpdateEntryID(context.Background(), v.JobId, entryId)
	}

	// 启动任务
	deployed.Cron.Start()
	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore start success.")
	// 停止任务
	defer deployed.Cron.Stop()
	select {}
}

// 添加任务 AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(job Job) (int, error) {
	id, err := deployed.Cron.AddJob(job.Expression(), job)
	return int(id), err
}

// 移除任务
func Remove(entryID int) chan struct{} {
	ch := make(chan struct{})
	go func() {
		deployed.Cron.Remove(cron.EntryID(entryID))
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Remove success ,info entryID :", entryID)
		ch <- struct{}{}
	}()
	return ch
}
