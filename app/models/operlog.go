package models

import (
	"time"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

// sys_operlog
type SysOperLog struct {
	OperId        int       `json:"operId" gorm:"primary_key;AUTO_INCREMENT"` // 日志编码
	Title         string    `json:"title" gorm:"size:255;"`                   // 操作模块
	BusinessType  string    `json:"businessType" gorm:"size:128;"`            // 操作类型
	BusinessTypes string    `json:"businessTypes" gorm:"size:128;"`
	Method        string    `json:"method" gorm:"size:128;"`         // 函数
	RequestMethod string    `json:"requestMethod" gorm:"size:128;"`  // 请求方式
	OperatorType  string    `json:"operatorType" gorm:"size:128;"`   // 操作类型
	OperName      string    `json:"operName" gorm:"size:128;"`       // 操作者
	DeptName      string    `json:"deptName" gorm:"size:128;"`       // 部门名称
	OperUrl       string    `json:"operUrl" gorm:"size:255;"`        // 访问地址
	OperIp        string    `json:"operIp" gorm:"size:128;"`         // 客户端ip
	OperLocation  string    `json:"operLocation" gorm:"size:128;"`   // 访问位置
	OperParam     string    `json:"operParam" gorm:"size:255;"`      // 请求参数
	Status        string    `json:"status" gorm:"size:4;"`           // 操作状态
	OperTime      time.Time `json:"operTime" gorm:"type:timestamp;"` // 操作时间
	JsonResult    string    `json:"jsonResult" gorm:"size:255;"`     // 返回数据
	CreateBy      string    `json:"createBy" gorm:"size:128;"`       // 创建人
	UpdateBy      string    `json:"updateBy" gorm:"size:128;"`       // 更新者
	DataScope     string    `json:"dataScope" gorm:"-"`              // 数据
	Params        string    `json:"params" gorm:"-"`                 // 参数
	Remark        string    `json:"remark" gorm:"size:255;"`         // 备注
	LatencyTime   string    `json:"latencyime" gorm:"size:128;"`     // 耗时
	UserAgent     string    `json:"userAgent" gorm:"size:255;"`      // ua
	Model
}

func (SysOperLog) TableName() string {
	return "sys_operlog"
}

func (e *SysOperLog) Get() (SysOperLog, error) {
	var doc SysOperLog

	table := deployed.DB.Table(e.TableName())
	if e.OperIp != "" {
		table = table.Where("oper_ip = ?", e.OperIp)
	}
	if e.OperId != 0 {
		table = table.Where("oper_id = ?", e.OperId)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *SysOperLog) GetPage(param paginator.Param) ([]SysOperLog, paginator.Info, error) {
	var doc []SysOperLog

	table := deployed.DB.Table(e.TableName())
	if e.OperIp != "" {
		table = table.Where("oper_ip = ?", e.OperIp)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}
	if e.OperName != "" {
		table = table.Where("oper_name = ?", e.OperName)
	}
	if e.BusinessType != "" {
		table = table.Where("business_type = ?", e.BusinessType)
	}

	ifc, err := iorm.QueryPages(table, param, &doc)
	if err != nil {
		return nil, ifc, err
	}
	return doc, ifc, nil
}

func (e *SysOperLog) Create() (SysOperLog, error) {
	var doc SysOperLog
	e.CreateBy = "0"
	e.UpdateBy = "0"
	result := deployed.DB.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *SysOperLog) Update(id int) (update SysOperLog, err error) {
	if err = deployed.DB.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	if err = deployed.DB.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *SysOperLog) BatchDelete(id []int) (Result bool, err error) {
	if err = deployed.DB.Table(e.TableName()).Where(" oper_id in (?)", id).Delete(&SysOperLog{}).Error; err != nil {
		return
	}
	Result = true
	return
}
