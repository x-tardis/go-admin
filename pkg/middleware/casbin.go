package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewCasbinEnforcer new casbin enforcer
func NewCasbinEnforcer(m model.Model, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
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

// NewCasbinEnforcerFromString new casbin enforcer from text
func NewCasbinEnforcerFromString(modelText string, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	md, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}
	return NewCasbinEnforcer(md, db)
}

// NedwCasbinEnforcerFromFile new casbin enforcer from file
func NewCasbinEnforcerFromFile(modelPath string, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	md, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}
	return NewCasbinEnforcer(md, db)
}
