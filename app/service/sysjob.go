package service

import (
	"context"
	"time"

	"github.com/thinkgos/http-middlewares/requestid"

	"github.com/x-tardis/go-admin/app/jobs"
	"github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type SysJob struct {
	Service
}

// RemoveJob 删除job
func (e *SysJob) RemoveJob(ctx context.Context, c *dto.GeneralDelDto) error {
	msgID := requestid.FromRequestID(ctx)
	item, err := models.CJob.Get(ctx, c.Id)
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}

	cn := jobs.Remove(item.EntryId)
	select {
	case <-cn:
		err = models.CJob.RemoveEntryID(ctx, item.EntryId)
		if err != nil {
			izap.Sugar.Errorf("msgID[%s] db error:%s", msgID, err)
		}
		return err
	case <-time.After(time.Second * 1):
		e.Msg = prompt.OperationTimeout.String()
	}
	return nil
}

// StartJob 启动任务
func (e *SysJob) StartJob(ctx context.Context, c *dto.GeneralGetDto) error {
	msgID := requestid.FromRequestID(ctx)

	item, err := models.CJob.Get(ctx, c.Id)
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}

	var job jobs.Job
	if item.JobType == 1 {
		job = &jobs.HttpJob{
			Base: jobs.Base{
				InvokeTarget:   item.InvokeTarget,
				Name:           item.JobName,
				JobId:          item.JobId,
				EntryId:        0,
				CronExpression: item.CronExpression,
				Args:           "",
			},
		}
	} else {
		job = &jobs.ExecJob{
			Base: jobs.Base{
				InvokeTarget:   item.InvokeTarget,
				Name:           item.JobName,
				JobId:          item.JobId,
				EntryId:        0,
				CronExpression: item.CronExpression,
				Args:           item.Args,
			},
		}
	}
	item.EntryId, err = jobs.AddJob(job)
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] jobs AddJob[HttpJob] error:%s", msgID, err)
		return err
	}

	err = models.CJob.UpdateEntryID(ctx, c.Id, item.EntryId)
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] db error:%s", msgID, err)
	}
	return err
}
