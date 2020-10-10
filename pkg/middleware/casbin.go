package middleware

import (
	"strings"

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

// NewCasbinEnforcerFromFile new casbin enforcer from file
func NewCasbinEnforcerFromFile(modelPath string, db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	md, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}
	return NewCasbinEnforcer(md, db)
}

// ParamMatch 自定义规则
// "/foo/bar","/foo" 匹配 "/foo/*"
func ParamMatch(key1 string, key2 string) bool {
	pos := strings.Index(key2, "*")
	if pos == -1 {
		return key1 == key2
	}

	if len(key1) >= pos {
		return key1[:pos] == key2[:pos]
	}

	return key1 == strings.TrimSuffix(key2[0:pos], "/")
}

// ParamMatchFunc wrap ParamMatch
func ParamMatchFunc(args ...interface{}) (interface{}, error) {
	return ParamMatch(args[0].(string), args[1].(string)), nil
}
