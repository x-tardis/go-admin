package models

import (
	"context"

	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
)

// CasbinRule casbin rule
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	PType string `gorm:"size:40;uniqueIndex:unique_index"` // type
	V0    string `gorm:"size:40;uniqueIndex:unique_index"` // role key
	V1    string `gorm:"size:40;uniqueIndex:unique_index"` // path
	V2    string `gorm:"size:40;uniqueIndex:unique_index"` // method
	V3    string `gorm:"size:40;uniqueIndex:unique_index"`
	V4    string `gorm:"size:40;uniqueIndex:unique_index"`
	V5    string `gorm:"size:40;uniqueIndex:unique_index"`
}

//TableName implement schema.Tabler interface
func (CasbinRule) TableName() string {
	return "sys_casbin_rule"
}

// CasbinRuleDB casbin rule db scope
func CasbinRuleDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(CasbinRule{})
	}
}

type cCasbinRule struct{}

// CCasbinRule 实例
var CCasbinRule = cCasbinRule{}

// BatchCreate 批量创建
func (cCasbinRule) BatchCreate(ctx context.Context, item []CasbinRule) ([]CasbinRule, error) {
	err := dao.DB.Scopes(CasbinRuleDB(ctx)).Create(&item).Error
	return item, err
}

// DeleteWithRoleName 通过角色名删除
func (cCasbinRule) DeleteWithRoleName(ctx context.Context, roleKey string) error {
	return dao.DB.Scopes(CasbinRuleDB(ctx)).
		Delete(&CasbinRule{}, "v0=?", roleKey).Error
}
