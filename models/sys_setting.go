package models

import (
	"context"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/trans"
)

// Setting 设置
type Setting struct {
	ID   int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // 主键
	Name string `json:"name" gorm:"type:varchar(256);"`       // 名称
	Logo string `json:"logo" gorm:"type:varchar(256);"`       // Logo
	Model
}

// TableName implement gorm.Tabler interface
func (Setting) TableName() string {
	return "sys_setting"
}

// SettingDB setting db scopes
func SettingDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Setting{})
	}
}

// UpSetting 更新设置
type UpSetting struct {
	Name string `json:"name" binding:"required"` // 名称
	Logo string `json:"logo" binding:"required"` // 头像
}

type cSetting struct{}

// CSetting 实例
var CSetting = new(cSetting)

// Get 获取
func (cSetting) Get(ctx context.Context) (item Setting, err error) {
	err = dao.DB.Scopes(SettingDB(ctx)).First(&item).Error
	return
}

// Update 修改
func (cSetting) Update(ctx context.Context, up UpSetting) (item Setting, err error) {
	err = dao.DB.Scopes(SettingDB(ctx)).
		Where("id=?", 1).
		Model(&item).Updates(&Setting{
		Logo: up.Logo,
		Name: up.Name,
	}).Error
	return
}
