package service

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/app/models"
	dto2 "github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/pkg/izap"
)

type SysJobSearch struct {
	paginator.Param `search:"-"`
	JobId           int    `form:"jobId" search:"type:exact;column:job_id;table:sys_job"`
	JobName         string `form:"jobName" search:"type:icontains;column:job_name;table:sys_job"`
	JobGroup        string `form:"jobGroup" search:"type:exact;column:job_group;table:sys_job"`
	CronExpression  string `form:"cronExpression" search:"type:exact;column:cron_expression;table:sys_job"`
	InvokeTarget    string `form:"invokeTarget" search:"type:exact;column:invoke_target;table:sys_job"`
	Status          int    `form:"status" search:"type:exact;column:status;table:sys_job"`
}

func (m *SysJobSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysJobSearch) Bind(ctx *gin.Context) error {
	err := ctx.ShouldBind(m)
	if err != nil {
		izap.Sugar.Errorf("MsgID[%s] Bind error: %s", err)
	}
	return err
}

func (m *SysJobSearch) Generate() dto2.Index {
	o := *m
	return &o
}

func (m *SysJobSearch) GetPaginatorParam() paginator.Param {
	m.Param.Inspect()
	return m.Param
}

type SysJobControl struct {
	JobId          uint   `json:"jobId"`
	JobName        string `json:"jobName" validate:"required"` // 名称
	JobGroup       string `json:"jobGroup"`                    // 任务分组
	JobType        int    `json:"jobType"`                     // 任务类型
	CronExpression string `json:"cronExpression"`              // cron表达式
	InvokeTarget   string `json:"invokeTarget"`                // 调用目标
	Args           string `json:"args"`                        // 目标参数
	MisfirePolicy  int    `json:"misfirePolicy"`               // 执行策略
	Concurrent     int    `json:"concurrent"`                  // 是否并发
	Status         int    `json:"status"`                      // 状态
	EntryId        int    `json:"entryId"`                     // job启动时返回的id
}

func (s *SysJobControl) Bind(ctx *gin.Context) error {
	return ctx.ShouldBind(s)
}

func (s *SysJobControl) Generate() dto2.Control {
	cp := *s
	return &cp
}

func (s *SysJobControl) GenerateM() (dto2.ActiveRecord, error) {
	return &models.Job{
		JobId:          s.JobId,
		JobName:        s.JobName,
		JobGroup:       s.JobGroup,
		JobType:        s.JobType,
		CronExpression: s.CronExpression,
		InvokeTarget:   s.InvokeTarget,
		Args:           s.Args,
		MisfirePolicy:  s.MisfirePolicy,
		Concurrent:     s.Concurrent,
		Status:         s.Status,
		EntryId:        s.EntryId,
	}, nil
}

func (s *SysJobControl) GetId() interface{} {
	return s.JobId
}

type SysJobById struct {
	dto2.ObjectById
}

func (s *SysJobById) Generate() dto2.Control {
	cp := *s
	return &cp
}

func (s *SysJobById) GenerateM() (dto2.ActiveRecord, error) {
	return &models.Job{}, nil
}

type SysJobItem struct {
	JobId          uint   `json:"jobId"`
	JobName        string `json:"jobName" validate:"required"` // 名称
	JobGroup       string `json:"jobGroup"`                    // 任务分组
	JobType        int    `json:"jobType"`                     // 任务类型
	CronExpression string `json:"cronExpression"`              // cron表达式
	InvokeTarget   string `json:"invokeTarget"`                // 调用目标
	Args           string `json:"args"`                        // 目标参数
	MisfirePolicy  int    `json:"misfirePolicy"`               // 执行策略
	Concurrent     int    `json:"concurrent"`                  // 是否并发
	Status         int    `json:"status"`                      // 状态
	EntryId        int    `json:"entryId"`                     // job启动时返回的id
}
