package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewCasbinSyncedEnforcer new casbin enforcer
func NewCasbinSyncedEnforcer(m model.Model, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	adapter, err := gormAdapter.NewAdapterByDBUsePrefix(db, "sys_")
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

// NewCasbinSyncedEnforcerFromString new casbin enforcer from text
func NewCasbinSyncedEnforcerFromString(modelText string, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	md, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}
	return NewCasbinSyncedEnforcer(md, db)
}

// NewCasbinSyncedEnforcerFromFile new casbin enforcer from file
func NewCasbinSyncedEnforcerFromFile(modelPath string, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	md, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}
	return NewCasbinSyncedEnforcer(md, db)
}
