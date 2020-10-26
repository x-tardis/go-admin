package models

import (
	"context"
	"strconv"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/trans"
)

type Job struct {
	JobId          uint   `json:"jobId" gorm:"primary_key;AUTO_INCREMENT"` // 编码
	JobName        string `json:"jobName" gorm:"size:255;"`                // 名称
	JobGroup       string `json:"jobGroup" gorm:"size:255;"`               // 任务分组
	JobType        int    `json:"jobType" gorm:"size:1;"`                  // 任务类型
	CronExpression string `json:"cronExpression" gorm:"size:255;"`         // cron表达式
	InvokeTarget   string `json:"invokeTarget" gorm:"size:255;"`           // 调用目标
	Args           string `json:"args" gorm:"size:255;"`                   // 目标参数
	MisfirePolicy  int    `json:"misfirePolicy" gorm:"size:255;"`          // 执行策略
	Concurrent     int    `json:"concurrent" gorm:"size:1;"`               // 是否并发
	Status         int    `json:"status" gorm:"size:1;"`                   // 状态
	EntryId        int    `json:"entry_id" gorm:"size:11;"`                // job启动时返回的id
	Creator        string `json:"creator" gorm:"size:128;"`                //
	Updator        string `json:"updator" gorm:"size:128;"`                //
	Model

	DataScope string `json:"dataScope" gorm:"-"`
}

func (Job) TableName() string {
	return "sys_job"
}

func JobDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Job{})
	}
}

type cJob struct{}

var CJob = cJob{}

func (e *Job) Generate() dto.ActiveRecord {
	o := *e
	return &o
}

func (e *Job) GetId() uint {
	return e.JobId
}

func (e *Job) SetCreator(createBy uint) {
	e.Creator = strconv.Itoa(int(createBy))
}

func (e *Job) SetUpdator(updateBy uint) {
	e.Updator = strconv.Itoa(int(updateBy))
}

// 获取SysJob带分页
func (cJob) QueryPage(ctx context.Context, param paginator.Param, v interface{}) ([]Job, paginator.Info, error) {
	var items []Job

	db := dao.DB.Scopes(JobDB(ctx), dto.MakeCondition(v))

	db = db.Scopes(DataScope(Job{}, jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return items, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, param, items)
	return items, info, err
}

func (cJob) Query(ctx context.Context) (items []Job, err error) {
	err = dao.DB.Scopes(JobDB(ctx)).
		Where("status=?", 2).Find(&items).Error
	return
}

// 获取SysJob
func (cJob) Get(ctx context.Context, id uint) (item Job, err error) {
	err = dao.DB.Scopes(JobDB(ctx)).
		Where("job_id=?", id).First(&item).Error
	return
}

// 创建SysJob
func (cJob) Create(ctx context.Context, item Job) (Job, error) {
	err := dao.DB.Scopes(JobDB(ctx)).Create(&item).Error
	return item, err
}

// 更新SysJob
func (cJob) Update(ctx context.Context, id uint, up Job) error {
	return dao.DB.Scopes(JobDB(ctx)).Where("job_id=?", id).Updates(&up).Error
}

func (cJob) RemoveAllEntryID(ctx context.Context) error {
	return dao.DB.Scopes(JobDB(ctx)).
		Where("entry_id > ?", 0).
		Update("entry_id", 0).Error
}

func (cJob) RemoveEntryID(ctx context.Context, entryID int) (err error) {
	return dao.DB.Scopes(JobDB(ctx)).
		Where("entry_id = ?", entryID).
		Update("entry_id", 0).Error
}

// 删除SysJob
func (cJob) Delete(ctx context.Context, id int) (err error) {
	return dao.DB.Scopes(JobDB(ctx)).
		Where("job_id=?", id).Delete(&Job{}).Error
}

// 批量删除
func (cJob) BatchDelete(ctx context.Context, ids []int) error {
	return dao.DB.Scopes(JobDB(ctx)).
		Where("job_id in ( ? )", ids).Delete(&Job{}).Error
}
