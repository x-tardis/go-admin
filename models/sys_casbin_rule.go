package models

import (
	"context"

	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
)

// CasbinRule casbin rule
type CasbinRule struct {
	PType string `json:"p_type" gorm:"size:100;"` // type
	V0    string `json:"v0" gorm:"size:100;"`     // role key
	V1    string `json:"v1" gorm:"size:100;"`     // path
	V2    string `json:"v2" gorm:"size:100;"`     // method
	V3    string `json:"v3" gorm:"size:100;"`
	V4    string `json:"v4" gorm:"size:100;"`
	V5    string `json:"v5" gorm:"size:100;"`
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
		Where("v0=?", roleKey).Delete(&CasbinRule{}).Error
}
