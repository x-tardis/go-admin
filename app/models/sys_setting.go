package models

import (
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

type SysSetting struct {
	ID   int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // 主键
	Name string `json:"name" gorm:"type:varchar(256);"`       // 名称
	Logo string `json:"logo" gorm:"type:varchar(256);"`       // Logo
	Model
}

func (SysSetting) TableName() string {
	return "sys_setting"
}

func SysSettingDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(SysSetting{})
	}
}

type UpSysSetting struct {
	Name string `json:"name" binding:"required"` // 名称
	Logo string `json:"logo" binding:"required"` // 头像
}

type CallSysSetting struct{}

// 查询
func (CallSysSetting) Get() (item SysSetting, err error) {
	err = deployed.DB.Scopes(SysSettingDB()).First(&item).Error
	return
}

// 修改
func (CallSysSetting) Update(up SysSetting) (item SysSetting, err error) {
	err = deployed.DB.Scopes(SysSettingDB()).Where("id=?", 1).Model(&item).Updates(&up).Error
	return
}
