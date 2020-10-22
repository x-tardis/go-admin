package models

import (
	"strconv"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/app/models/dao"
	dto2 "github.com/x-tardis/go-admin/app/service/dto"
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

func JobDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(Job{})
	}
}

type cJob struct{}

var CJob = new(cJob)

func (e *Job) Generate() dto2.ActiveRecord {
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
func (e *Job) GetPage(pageSize int, pageIndex int, v interface{}, list interface{}) (int, error) {
	table := dao.DB.Table(e.TableName()).Scopes(dto2.MakeCondition(v))

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	// dataPermission := new(DataPermission)
	userid, _ := strconv.Atoi(e.DataScope)

	var count int64

	if err := table.Scopes(DataScopes(e.TableName(), userid),
		iorm.Paginate(paginator.Param{pageIndex, pageSize})).
		Find(list).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (cJob) Query() (items []Job, err error) {
	err = dao.DB.Scopes(JobDB()).
		Where("status=?", 2).Find(&items).Error
	return
}

// 获取SysJob
func (cJob) Get(id uint) (item Job, err error) {
	err = dao.DB.Scopes(JobDB()).
		Where("job_id=?", id).First(&item).Error
	return
}

// 创建SysJob
func (cJob) Create(item Job) (Job, error) {
	err := dao.DB.Scopes(JobDB()).Create(&item).Error
	return item, err
}

// 更新SysJob
func (cJob) Update(id uint, e Job) error {
	return dao.DB.Scopes(JobDB()).Where("job_id=?", id).Updates(&e).Error
}

func (cJob) RemoveAllEntryID() error {
	return dao.DB.Scopes(JobDB()).
		Where("entry_id > ?", 0).
		Update("entry_id", 0).Error
}

func (cJob) RemoveEntryID(entryID int) (err error) {
	return dao.DB.Scopes(JobDB()).
		Where("entry_id = ?", entryID).
		Update("entry_id", 0).Error
}

// 删除SysJob
func (cJob) Delete(id int) (err error) {
	return dao.DB.Scopes(JobDB()).
		Where("job_id=?", id).Delete(&Job{}).Error
}

// 批量删除
func (cJob) BatchDelete(ids []int) error {
	return dao.DB.Scopes(JobDB()).
		Where("job_id in ( ? )", ids).Delete(&Job{}).Error
}
