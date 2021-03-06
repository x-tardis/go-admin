package models

import (
	"context"
	"time"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
)

// LoginLog 登录记录
type LoginLog struct {
	InfoId    int       `json:"infoId" gorm:"primary_key;auto_increment;"` // 主键
	Username  string    `json:"username" gorm:"size:128;"`                 // 用户名
	Status    string    `json:"status" gorm:"size:4;"`                     // 登录状态
	Ip        string    `json:"ip" gorm:"size:255;"`                       // 登录ip地址
	Location  string    `json:"location" gorm:"size:255;"`                 // 登录ip归属地
	Browser   string    `json:"browser" gorm:"size:255;"`                  // 浏览器
	Os        string    `json:"os" gorm:"size:255;"`                       // 操作系统
	Platform  string    `json:"platform" gorm:"size:255;"`                 // 系统平台
	LoginTime time.Time `json:"loginTime" gorm:"type:timestamp;"`          // 登录时间
	Remark    string    `json:"remark" gorm:"size:255;"`                   // 备注
	Msg       string    `json:"msg" gorm:"size:255;"`                      // 登录信息
	Creator   string    `json:"creator" gorm:"size:128;"`                  // 创建人
	Updator   string    `json:"updator" gorm:"size:128;"`                  // 更新者
	Model

	DataScope string `json:"dataScope" gorm:"-"` // 数据
	Params    string `json:"params" gorm:"-"`    //
}

// TableName implement schema.Tabler interface
func (LoginLog) TableName() string {
	return "sys_loginlog"
}

// LoginLogDB login log db scopes
func LoginLogDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(LoginLog{})
	}
}

// LoginLogQueryParam 查询参数
type LoginLogQueryParam struct {
	Username string `form:"username"` // 用户名
	Ip       string `form:"ip"`       // ip地址
	Status   string `form:"status"`   // 状态
	paginator.Param
}

type cLoginLog struct{}

// CLoginLog 实例
var CLoginLog = cLoginLog{}

// QueryPage 查询,分页
func (cLoginLog) QueryPage(ctx context.Context, qp LoginLogQueryParam) ([]LoginLog, paginator.Info, error) {
	var items []LoginLog

	db := dao.DB.Scopes(LoginLogDB(ctx))
	if qp.Username != "" {
		db = db.Where("username=?", qp.Username)
	}
	if qp.Ip != "" {
		db = db.Where("ip=?", qp.Ip)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}
	db = db.Order("info_id desc")

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// Get 获取
func (cLoginLog) Get(ctx context.Context, id int) (item LoginLog, err error) {
	err = dao.DB.Scopes(LoginLogDB(ctx)).
		First(&item, "info_id=?", id).Error
	return
}

// Create 创建
func (cLoginLog) Create(ctx context.Context, item LoginLog) (LoginLog, error) {
	item.Creator = "0"
	item.Updator = "0"
	err := dao.DB.Scopes(LoginLogDB(ctx)).Create(&item).Error
	return item, err
}

// Update 更新
func (sf cLoginLog) Update(ctx context.Context, id int, up LoginLog) error {
	item, err := sf.Get(ctx, id)
	if err != nil {
		return err
	}
	return dao.DB.Scopes(LoginLogDB(ctx)).Model(&item).Updates(&up).Error
}

// Delete 删除
func (cLoginLog) Delete(ctx context.Context, id int) error {
	return dao.DB.Scopes(LoginLogDB(ctx)).
		Delete(&LoginLog{}, "info_id=?", id).Error
}

// BatchDelete 批量删除id
func (sf cLoginLog) BatchDelete(ctx context.Context, ids []int) error {
	switch len(ids) {
	case 0:
		return nil
	case 1:
		return sf.Delete(ctx, ids[0])
	default:
		return dao.DB.Scopes(LoginLogDB(ctx)).
			Delete(&LoginLog{}, "info_id in (?)", ids).Error
	}
}

// Clean 清空日志
func (cLoginLog) Clean(ctx context.Context) error {
	return dao.DB.Scopes(LoginLogDB(ctx)).Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&LoginLog{}).Error
}
