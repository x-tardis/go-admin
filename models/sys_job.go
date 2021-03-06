package models

import (
	"context"
	"strconv"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// Job job
type Job struct {
	JobId          uint   `json:"jobId" gorm:"primary_key;AUTO_INCREMENT"` // 主键
	JobName        string `json:"jobName" gorm:"size:255;"`                // 名称
	JobGroup       string `json:"jobGroup" gorm:"size:255;"`               // 分组
	JobType        int    `json:"jobType" gorm:"size:1;"`                  // 类型
	CronExpression string `json:"cronExpression" gorm:"size:255;"`         // cron表达式
	InvokeTarget   string `json:"invokeTarget" gorm:"size:255;"`           // 调用目标
	Args           string `json:"args" gorm:"size:255;"`                   // 目标参数
	MisfirePolicy  int    `json:"misfirePolicy" gorm:"size:255;"`          // 执行策略
	Concurrent     int    `json:"concurrent" gorm:"size:1;"`               // 是否并发
	Status         int    `json:"status" gorm:"size:1;"`                   // 状态
	EntryId        int    `json:"entry_id" gorm:"size:11;"`                // job启动时cron返回的条目id
	Creator        string `json:"creator" gorm:"size:128;"`                // 创建者
	Updator        string `json:"updator" gorm:"size:128;"`                // 更新者
	Model

	DataScope string `json:"dataScope" gorm:"-"`
}

// TableName implement schema.Tabler interface
func (Job) TableName() string {
	return "sys_job"
}

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

// JobDB job db scope
func JobDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Job{})
	}
}

type cJob struct{}

// 实例
var CJob = cJob{}

// QueryPage 查询,分页
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

// Query 查询
func (cJob) Query(ctx context.Context) (items []Job, err error) {
	err = dao.DB.Scopes(JobDB(ctx)).
		Find(&items, "status=?", 2).Error
	return
}

// Get 获取
func (cJob) Get(ctx context.Context, id uint) (item Job, err error) {
	err = dao.DB.Scopes(JobDB(ctx)).
		First(&item, "job_id=?", id).Error
	return
}

// Create 创建
func (cJob) Create(ctx context.Context, item Job) (Job, error) {
	err := dao.DB.Scopes(JobDB(ctx)).Create(&item).Error
	return item, err
}

// Update 更新SysJob
func (cJob) Update(ctx context.Context, id uint, up Job) error {
	return dao.DB.Scopes(JobDB(ctx)).Where("job_id=?", id).Updates(&up).Error
}

func (cJob) UpdateEntryID(ctx context.Context, id uint, entryID int) error {
	return dao.DB.Scopes(JobDB(ctx)).
		Where("job_id=?", id).Update("entry_id", entryID).Error
}

// Delete 删除
func (cJob) Delete(ctx context.Context, id int) (err error) {
	return dao.DB.Scopes(JobDB(ctx)).
		Delete(&Job{}, "job_id=?", id).Error
}

// BatchDelete 批量删除
func (cJob) BatchDelete(ctx context.Context, ids []int) error {
	return dao.DB.Scopes(JobDB(ctx)).
		Delete(&Job{}, "job_id in ( ? )", ids).Error
}

// RemoveAllEntryID 删除所有entry id
func (cJob) RemoveAllEntryID(ctx context.Context) error {
	return dao.DB.Scopes(JobDB(ctx)).
		Where("entry_id > ?", 0).Update("entry_id", 0).Error
}

// RemoveEntryID 删除指定entry id
func (cJob) RemoveEntryID(ctx context.Context, entryID int) error {
	return dao.DB.Scopes(JobDB(ctx)).
		Where("entry_id=?", entryID).Update("entry_id", 0).Error
}
