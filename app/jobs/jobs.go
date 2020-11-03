package jobs

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/models"
)

var jobList map[string]JobExec

// 需要将定义的struct 添加到字典中；
// 字典 key 可以配置到 自动任务 调用目标 中；
func initJob() {
	jobList = map[string]JobExec{
		"ExamplesOne": ExamplesOne{},
		// ...
	}
}

// 初始化
func Startup() {
	initJob()
	log.Println("[INFO] JobCore Starting...")

	deployed.Cron = deployed.NewCron()

	jobItems, err := models.CJob.Query(context.Background())
	if err != nil {
		log.Println("[ERROR] JobCore init error", err)
	}

	log.Println("[INFO] JobCore total:0", len(jobItems))

	err = models.CJob.RemoveAllEntryID(context.Background())
	if err != nil {
		log.Println("[ERROR] JobCore remove all entry id failed", err)
	}

	for _, v := range jobItems {
		job := Convert(v)
		entryId, err := AddJob(job)
		if err != nil {
			continue
		}
		models.CJob.UpdateEntryID(context.Background(), v.JobId, entryId) // nolint: errcheck
	}

	// 启动任务
	deployed.Cron.Start()
	log.Println("[INFO] JobCore start success.")
}

func Stop() {
	deployed.Cron.Stop()
}

// 添加任务 AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(job Job) (int, error) {
	id, err := deployed.Cron.AddJob(job.Expression(), job)
	job.SetEntryId(int(id))
	return int(id), err
}

// 移除任务
func Remove(entryID int) chan struct{} {
	ch := make(chan struct{})
	go func() {
		deployed.Cron.Remove(cron.EntryID(entryID))
		log.Println("[INFO] JobCore Remove success ,info entryID :", entryID)
		ch <- struct{}{}
	}()
	return ch
}

func Convert(item models.Job) Job {
	base := Base{
		item.JobId,
		0,
		item.JobName,
		item.InvokeTarget,
		item.CronExpression,
		item.Args,
	}
	switch item.JobType {
	case 1:
		return &HttpJob{base}
	case 2:
		return &ExecJob{base}
	}
	return nil
}
