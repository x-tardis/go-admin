package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewCasbinSyncedEnforcer new casbin enforcer with table name sys_casbin_rule
func NewCasbinSyncedEnforcer(m model.Model, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	adapter, err := gormAdapter.NewAdapterByDBUseTableName(db, "sys", "casbin_rule")
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	err = e.LoadPolicy()
	return e, err
}

// NewCasbinSyncedEnforcerFromString new casbin enforcer from text with table name sys_casbin_rule
func NewCasbinSyncedEnforcerFromString(modelText string, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	md, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}
	return NewCasbinSyncedEnforcer(md, db)
}

// NewCasbinSyncedEnforcerFromFile new casbin enforcer from file  with table name sys_casbin_rule
func NewCasbinSyncedEnforcerFromFile(modelPath string, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	md, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}
	return NewCasbinSyncedEnforcer(md, db)
}

// NewCasbinEnforcer new casbin enforcer with table name sys_casbin_rule
func NewCasbinEnforcer(m model.Model, db *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := gormAdapter.NewAdapterByDBUseTableName(db, "sys", "casbin_rule")
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	err = e.LoadPolicy()
	return e, err
}

// NewCasbinEnforcerFromString new casbin enforcer from text with table name sys_casbin_rule
func NewCasbinEnforcerFromString(modelText string, db *gorm.DB) (*casbin.Enforcer, error) {
	md, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}
	return NewCasbinEnforcer(md, db)
}

// NewCasbinEnforcerFromFile new casbin enforcer from file with table name sys_casbin_rule
func NewCasbinEnforcerFromFile(modelPath string, db *gorm.DB) (*casbin.Enforcer, error) {
	md, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}
	return NewCasbinEnforcer(md, db)
}
