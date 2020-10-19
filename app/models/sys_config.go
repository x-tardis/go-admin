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

type Config struct {
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

func (Config) TableName() string {
	return "sys_config"
}

func ConfigDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(Config{})
	}
}

// ConfigQueryParam 查询参数
type ConfigQueryParam struct {
	ConfigName string `form:"configName"`
	ConfigKey  string `form:"configKey"`
	ConfigType string `form:"configType"`
	paginator.Param
}

type cConfig struct{}

var CConfig = new(cConfig)

// 获取 Config
func (cConfig) GetWithKey(_ context.Context, key string) (item Config, err error) {
	err = deployed.DB.Scopes(ConfigDB()).
		Where("config_key = ?", key).First(&item).Error
	return
}

// 获取 Config
func (cConfig) Get(_ context.Context, id int) (item Config, err error) {
	err = deployed.DB.Scopes(ConfigDB()).
		Where("config_id = ?", id).First(&item).Error
	return
}

func (cConfig) QueryPage(ctx context.Context, qp ConfigQueryParam) ([]Config, paginator.Info, error) {
	var items []Config

	db := deployed.DB.Scopes(ConfigDB())

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

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// Config 创建
func (cConfig) Create(ctx context.Context, item Config) (Config, error) {
	var i int64

	item.Creator = jwtauth.FromUserIdStr(ctx)
	deployed.DB.Scopes(ConfigDB()).
		Where("config_name=? or config_key = ?", item.ConfigName, item.ConfigKey).
		Count(&i)
	if i > 0 {
		return item, iorm.ErrObjectAlreadyExist
	}

	result := deployed.DB.Scopes(ConfigDB()).Create(&item)
	if err := result.Error; err != nil {
		return item, err
	}
	return item, nil
}

func (cConfig) Update(ctx context.Context, id int, item Config) (update Config, err error) {
	item.Updator = jwtauth.FromUserIdStr(ctx)
	if err = deployed.DB.Scopes(ConfigDB()).
		Where("config_id = ?", id).First(&update).Error; err != nil {
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
	err = deployed.DB.Scopes(ConfigDB()).Model(&update).Updates(&item).Error
	return
}

func (cConfig) Delete(_ context.Context, id int) (err error) {
	return deployed.DB.Scopes(ConfigDB()).
		Where("config_id = ?", id).Delete(&Config{}).Error
}

func (cConfig) BatchDelete(_ context.Context, id []int) error {
	return deployed.DB.Scopes(ConfigDB()).
		Where("config_id in (?)", id).Delete(&Config{}).Error
}
