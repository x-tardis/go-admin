package service

import (
	"context"
	"time"

	"github.com/thinkgos/http-middlewares/requestid"

	"github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/jobs"
	"github.com/x-tardis/go-admin/pkg/servers/codes"
)

type Service struct {
	Msg string
}

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
		e.Msg = codes.OperationTimeout
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

	entryId, err := jobs.AddJob(jobs.Convert(item))
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] jobs AddJob[HttpJob] error:%s", msgID, err)
		return err
	}

	err = models.CJob.UpdateEntryID(ctx, c.Id, entryId)
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] db error:%s", msgID, err)
	}
	return err
}
