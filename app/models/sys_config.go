package models

import (
	"context"
	"errors"
	_ "time"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type SysConfig struct {
	ConfigId    int    `json:"configId" gorm:"primary_key;auto_increment;"` // 主键
	ConfigName  string `json:"configName" gorm:"size:128;"`                 // 参数名称
	ConfigKey   string `json:"configKey" gorm:"size:128;"`                  // 参数键名
	ConfigValue string `json:"configValue" gorm:"size:255;"`                // 参数键值
	ConfigType  string `json:"configType" gorm:"size:64;"`                  // 是否系统内置
	Remark      string `json:"remark" gorm:"size:128;"`                     // 备注
	Creator     string `json:"creator" gorm:"size:128;"`                    // 创建者
	Updator     string `json:"updator" gorm:"size:128;"`                    // 更新者
	Model

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

// SysConfigQueryParam 查询参数
type SysConfigQueryParam struct {
	ConfigName string `form:"configName"`
	ConfigKey  string `form:"configKey"`
	ConfigType string `form:"configType"`
	paginator.Param
}

func (SysConfig) TableName() string {
	return "sys_config"
}

func SysConfigDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(SysConfig{})
	}
}

type CallSysConfig struct{}

// 获取 Config
func (CallSysConfig) GetWithKey(_ context.Context, key string) (item SysConfig, err error) {
	err = deployed.DB.Scopes(SysConfigDB()).
		Where("config_key = ?", key).First(&item).Error
	return
}

// 获取 Config
func (CallSysConfig) Get(_ context.Context, id int) (item SysConfig, err error) {
	err = deployed.DB.Scopes(SysConfigDB()).
		Where("config_id = ?", id).First(&item).Error
	return
}

func (CallSysConfig) QueryPage(ctx context.Context, qp SysConfigQueryParam) ([]SysConfig, paginator.Info, error) {
	var items []SysConfig

	db := deployed.DB.Scopes(SysConfigDB())

	if qp.ConfigName != "" {
		db = db.Where("config_name=?", qp.ConfigName)
	}
	if qp.ConfigKey != "" {
		db = db.Where("config_key=?", qp.ConfigKey)
	}
	if qp.ConfigType != "" {
		db = db.Where("config_type=?", qp.ConfigType)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId = jwtauth.FromUserId(ctx)
	db, err := dataPermission.GetDataScope("sys_config", db)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	ifc, err := iorm.QueryPages(db, qp.Param, &items)
	if err != nil {
		return nil, ifc, err
	}
	return items, ifc, nil
}

// Config 创建
func (CallSysConfig) Create(ctx context.Context, item SysConfig) (SysConfig, error) {
	var i int64

	item.Creator = jwtauth.FromUserIdStr(ctx)
	deployed.DB.Scopes(SysConfigDB()).Where("config_name=? or config_key = ?", item.ConfigName, item.ConfigKey).Count(&i)
	if i > 0 {
		return item, iorm.ErrObjectAlreadyExist
	}

	result := deployed.DB.Scopes(SysConfigDB()).Create(&item)
	if err := result.Error; err != nil {
		return item, err
	}
	return item, nil
}

func (CallSysConfig) Update(ctx context.Context, id int, item SysConfig) (update SysConfig, err error) {
	item.Updator = jwtauth.FromUserIdStr(ctx)
	if err = deployed.DB.Scopes(SysConfigDB()).Where("config_id = ?", id).First(&update).Error; err != nil {
		return
	}

	if item.ConfigName != "" && item.ConfigName != update.ConfigName {
		return update, errors.New("参数名称不允许修改！")
	}

	if item.ConfigKey != "" && item.ConfigKey != update.ConfigKey {
		return update, errors.New("参数键名不允许修改！")
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	if err = deployed.DB.Scopes(SysConfigDB()).Model(&update).Updates(&item).Error; err != nil {
		return
	}
	return
}

func (CallSysConfig) Delete(_ context.Context, id int) (err error) {
	return deployed.DB.Scopes(SysConfigDB()).
		Where("config_id = ?", id).Delete(&SysConfig{}).Error
}

func (CallSysConfig) BatchDelete(_ context.Context, id []int) error {
	return deployed.DB.Scopes(SysConfigDB()).
		Where("config_id in (?)", id).Delete(&SysConfig{}).Error
}
