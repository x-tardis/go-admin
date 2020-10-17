package models

import (
	"context"
	"time"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

type LoginLog struct {
	InfoId        int       `json:"infoId" gorm:"primary_key;auto_increment;"` // 主键
	Username      string    `json:"username" gorm:"size:128;"`                 // 用户名
	Status        string    `json:"status" gorm:"size:4;"`                     // 状态
	Ipaddr        string    `json:"ipaddr" gorm:"size:255;"`                   // ip地址
	LoginLocation string    `json:"loginLocation" gorm:"size:255;"`            // 归属地
	Browser       string    `json:"browser" gorm:"size:255;"`                  // 浏览器
	Os            string    `json:"os" gorm:"size:255;"`                       // 系统
	Platform      string    `json:"platform" gorm:"size:255;"`                 // 固件
	LoginTime     time.Time `json:"loginTime" gorm:"type:timestamp;"`          // 登录时间
	Creator       string    `json:"creator" gorm:"size:128;"`                  // 创建人
	Updator       string    `json:"updator" gorm:"size:128;"`                  // 更新者
	Remark        string    `json:"remark" gorm:"size:255;"`                   // 备注
	Msg           string    `json:"msg" gorm:"size:255;"`
	Model

	DataScope string `json:"dataScope" gorm:"-"` // 数据
	Params    string `json:"params" gorm:"-"`    //
}

func (LoginLog) TableName() string {
	return "sys_loginlog"
}

func LoginLogDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(LoginLog{})
	}
}

type LoginLogQueryParam struct {
	Username string `form:"username"` // 用户名
	Status   string `form:"status"`   // 状态
	Ipaddr   string `form:"ipaddr"`   // ip地址
	paginator.Param
}

type CallLoginLog struct{}

func (CallLoginLog) Get(id int) (item LoginLog, err error) {
	err = deployed.DB.Scopes(LoginLogDB()).
		Where("info_id = ?", id).First(&item).Error
	return
}

func (CallLoginLog) QueryPage(_ context.Context, qp LoginLogQueryParam) ([]LoginLog, paginator.Info, error) {
	db := deployed.DB.Scopes(LoginLogDB())
	if qp.Ipaddr != "" {
		db = db.Where("ipaddr=?", qp.Ipaddr)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}
	if qp.Username != "" {
		db = db.Where("username=?", qp.Username)
	}
	db = db.Order("info_id desc")

	var items []LoginLog

	ifc, err := iorm.QueryPages(db, qp.Param, &items)
	if err != nil {
		return nil, ifc, err
	}
	return items, ifc, nil
}

func (CallLoginLog) Create(_ context.Context, item LoginLog) (LoginLog, error) {
	item.Creator = "0"
	item.Updator = "0"
	err := deployed.DB.Scopes(LoginLogDB()).Create(&item).Error
	return item, err
}

func (CallLoginLog) Update(_ context.Context, id int, up LoginLog) (item LoginLog, err error) {
	if err = deployed.DB.Scopes(LoginLogDB()).First(&item, id).Error; err != nil {
		return
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据,只修改非零的数据
	if err = deployed.DB.Scopes(LoginLogDB()).Model(&item).Updates(&up).Error; err != nil {
		return
	}
	return
}

// BatchDelete 批量删除id
func (CallLoginLog) BatchDelete(_ context.Context, id []int) error {
	return deployed.DB.Scopes(LoginLogDB()).
		Where("info_id in (?)", id).Delete(&LoginLog{}).Error
}

// Clean 清空日志
func (CallLoginLog) Clean(_ context.Context) error {
	return deployed.DB.Scopes(LoginLogDB()).Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&LoginLog{}).Error
}
