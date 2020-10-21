package service

import (
	"time"

	"github.com/x-tardis/go-admin/app/jobs"
	"github.com/x-tardis/go-admin/app/models"
	dto2 "github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type SysJob struct {
	Service
}

// RemoveJob 删除job
func (e *SysJob) RemoveJob(c *dto2.GeneralDelDto) error {
	var err error
	var data models.Job
	msgID := e.MsgID
	data.JobId = c.Id
	err = e.Orm.Table(data.TableName()).First(&data).Error
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	cn := jobs.Remove(data.EntryId)

	select {
	case res := <-cn:
		if res {
			err = e.Orm.Table(data.TableName()).Where("entry_id = ?", data.EntryId).Update("entry_id", 0).Error
			if err != nil {
				izap.Sugar.Errorf("msgID[%s] db error:%s", msgID, err)
			}
			return err
		}
	case <-time.After(time.Second * 1):
		e.Msg = prompt.OperationTimeout.String()
		return nil
	}
	return nil
}

// StartJob 启动任务
func (e *SysJob) StartJob(c *dto2.GeneralGetDto) error {
	var data models.Job
	var err error
	msgID := e.MsgID
	err = e.Orm.Table(data.TableName()).First(&data, c.Id).Error
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if data.JobType == 1 {
		var j = &jobs.HttpJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		data.EntryId, err = jobs.AddJob(j)
		if err != nil {
			izap.Sugar.Errorf("msgID[%s] jobs AddJob[HttpJob] error:%s", msgID, err)
		}
	} else {
		var j = &jobs.ExecJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		j.Args = data.Args
		data.EntryId, err = jobs.AddJob(j)
		if err != nil {
			izap.Sugar.Errorf("msgID[%s] jobs AddJob[ExecJob] error:%s", msgID, err)
		}
	}
	if err != nil {
		return err
	}

	err = e.Orm.Table(data.TableName()).Where(c.Id).Updates(&data).Error
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] db error:%s", msgID, err)
	}
	return err
}
