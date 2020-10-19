package models

import (
	"context"
	"time"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

// OperLog
type OperLog struct {
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
	Remark        string    `json:"remark" gorm:"size:255;"`         // 备注
	LatencyTime   string    `json:"latencyime" gorm:"size:128;"`     // 耗时
	UserAgent     string    `json:"userAgent" gorm:"size:255;"`      // ua
	Creator       string    `json:"creator" gorm:"size:128;"`        // 创建人
	Updator       string    `json:"updator" gorm:"size:128;"`        // 更新者
	Model

	DataScope string `json:"dataScope" gorm:"-"` // 数据
	Params    string `json:"params" gorm:"-"`    // 参数
}

func (OperLog) TableName() string {
	return "sys_operlog"
}

func OperLogDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(OperLog{})
	}
}

type OperLogQueryParam struct {
	Title        string `form:"title"`        // 操作模块
	OperName     string `form:"operName"`     // 操作人员
	OperIp       string `form:"operIp"`       // 客户端ip
	BusinessType string `form:"businessType"` // 操作类型
	Status       string `form:"status"`       // 操作状态
	paginator.Param
}

type cOperLog struct{}

var COperLog = new(cOperLog)

func (cOperLog) Get(_ context.Context, id int) (item OperLog, err error) {
	err = deployed.DB.Scopes(OperLogDB()).
		Where("oper_id=?", id).First(&item).Error
	return
}

func (cOperLog) QueryPage(_ context.Context, qp OperLogQueryParam) ([]OperLog, paginator.Info, error) {
	var items []OperLog

	db := deployed.DB.Scopes(OperLogDB())
	if qp.Title != "" {
		db = db.Where("title=?", qp.Title)
	}
	if qp.OperIp != "" {
		db = db.Where("oper_ip=?", qp.OperIp)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}
	if qp.OperName != "" {
		db = db.Where("oper_name=?", qp.OperName)
	}
	if qp.BusinessType != "" {
		db = db.Where("business_type=?", qp.BusinessType)
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

func (cOperLog) Create(_ context.Context, item OperLog) (OperLog, error) {
	item.Creator = "0"
	item.Updator = "0"
	err := deployed.DB.Scopes(OperLogDB()).Create(&item).Error
	return item, err
}

func (cOperLog) Update(id int, up OperLog) (item OperLog, err error) {
	if err = deployed.DB.Scopes(OperLogDB()).First(&item, id).Error; err != nil {
		return
	}
	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = deployed.DB.Scopes(OperLogDB()).Model(&item).Updates(&up).Error
	return
}

func (cOperLog) BatchDelete(_ context.Context, id []int) error {
	return deployed.DB.Scopes(OperLogDB()).
		Where(" oper_id in (?)", id).Delete(&OperLog{}).Error
}

func (cOperLog) Clean(_ context.Context) error {
	return deployed.DB.Scopes(OperLogDB()).Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&OperLog{}).Error
}
