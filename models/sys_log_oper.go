package models

import (
	"context"
	"time"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// OperLog operate log
type OperLog struct {
	OperId        int       `json:"operId" gorm:"primary_key;AUTO_INCREMENT"` // 主鍵
	Title         string    `json:"title" gorm:"size:255;"`                   // 操作模块
	BusinessType  string    `json:"businessType" gorm:"size:128;"`            // 操作类型
	Method        string    `json:"method" gorm:"size:128;"`                  // 函数
	RequestMethod string    `json:"requestMethod" gorm:"size:128;"`           // 请求方式
	OperatorType  string    `json:"operatorType" gorm:"size:128;"`            // 操作类型
	OperName      string    `json:"operName" gorm:"size:128;"`                // 操作者
	DeptName      string    `json:"deptName" gorm:"size:128;"`                // 部门名称
	OperUrl       string    `json:"operUrl" gorm:"size:255;"`                 // 访问地址
	OperIp        string    `json:"operIp" gorm:"size:128;"`                  // 客户端ip
	OperLocation  string    `json:"operLocation" gorm:"size:128;"`            // 访问位置
	OperParam     string    `json:"operParam" gorm:"size:255;"`               // 请求参数
	Status        string    `json:"status" gorm:"size:4;"`                    // 操作状态
	OperTime      time.Time `json:"operTime" gorm:"type:timestamp;"`          // 操作时间
	JsonResult    string    `json:"jsonResult" gorm:"size:255;"`              // 返回数据
	Remark        string    `json:"remark" gorm:"size:255;"`                  // 备注
	LatencyTime   string    `json:"latencyime" gorm:"size:128;"`              // 耗时
	UserAgent     string    `json:"userAgent" gorm:"size:255;"`               // user_agent
	Creator       string    `json:"creator" gorm:"size:128;"`                 // 创建人
	Updator       string    `json:"updator" gorm:"size:128;"`                 // 更新者
	Model

	DataScope string `json:"dataScope" gorm:"-"` // 数据
	Params    string `json:"params" gorm:"-"`    // 参数
}

// TableName implement schema.Tabler interface
func (OperLog) TableName() string {
	return "sys_operlog"
}

// OperLogDB openate log db scope
func OperLogDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(OperLog{})
	}
}

// OperLogQueryParam 查询参数
type OperLogQueryParam struct {
	Title        string `form:"title"`        // 操作模块
	OperName     string `form:"operName"`     // 操作人员
	OperIp       string `form:"operIp"`       // 客户端ip
	BusinessType string `form:"businessType"` // 操作类型
	Status       string `form:"status"`       // 操作状态
	paginator.Param
}

type cOperLog struct{}

// COperLog 实例
var COperLog = cOperLog{}

// QueryPage 查询,分页
func (cOperLog) QueryPage(ctx context.Context, qp OperLogQueryParam) ([]OperLog, paginator.Info, error) {
	var items []OperLog

	db := dao.DB.Scopes(OperLogDB(ctx))
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

// Get 获取
func (cOperLog) Get(ctx context.Context, id int) (item OperLog, err error) {
	err = dao.DB.Scopes(OperLogDB(ctx)).
		First(&item, "oper_id=?", id).Error
	return
}

// Create 创建
func (cOperLog) Create(ctx context.Context, item OperLog) (OperLog, error) {
	item.Creator = "0"
	item.Updator = "0"
	err := dao.DB.Scopes(OperLogDB(ctx)).Create(&item).Error
	return item, err
}

// Update 更新
func (sf cOperLog) Update(ctx context.Context, id int, up OperLog) error {
	item, err := sf.Get(ctx, id)
	if err != nil {
		return err
	}
	up.Updator = jwtauth.FromUserIdStr(ctx)
	return dao.DB.Scopes(OperLogDB(ctx)).Model(&item).Updates(&up).Error
}

// BatchDelete 批量删除
func (cOperLog) BatchDelete(ctx context.Context, id []int) error {
	return dao.DB.Scopes(OperLogDB(ctx)).
		Delete(&OperLog{}, "oper_id in (?)", id).Error
}

// Clean 清空
func (cOperLog) Clean(ctx context.Context) error {
	return dao.DB.Scopes(OperLogDB(ctx)).Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&OperLog{}).Error
}
