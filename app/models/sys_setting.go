package models

import (
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

type Setting struct {
	ID   int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // 主键
	Name string `json:"name" gorm:"type:varchar(256);"`       // 名称
	Logo string `json:"logo" gorm:"type:varchar(256);"`       // Logo
	Model
}

func (Setting) TableName() string {
	return "sys_setting"
}

func SettingDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(Setting{})
	}
}

type UpSetting struct {
	Name string `json:"name" binding:"required"` // 名称
	Logo string `json:"logo" binding:"required"` // 头像
}

type cSetting struct{}

var CSetting = new(cSetting)

// 查询
func (cSetting) Get() (item Setting, err error) {
	err = deployed.DB.Scopes(SettingDB()).First(&item).Error
	return
}

// 修改
func (cSetting) Update(up Setting) (item Setting, err error) {
	err = deployed.DB.Scopes(SettingDB()).Where("id=?", 1).Model(&item).Updates(&up).Error
	return
}
