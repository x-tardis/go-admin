package models

import (
	"context"
	"errors"
	_ "time"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/trans"
)

// Config参数配置
type Config struct {
	ConfigId    int    `json:"configId" gorm:"primary_key;auto_increment;"` // 主键
	ConfigName  string `json:"configName" gorm:"size:128;"`                 // 名称
	ConfigKey   string `json:"configKey" gorm:"size:128;"`                  // 键名
	ConfigValue string `json:"configValue" gorm:"size:255;"`                // 键值
	ConfigType  string `json:"configType" gorm:"size:64;"`                  // 是否系统内置
	Remark      string `json:"remark" gorm:"size:128;"`                     // 备注
	Creator     string `json:"creator" gorm:"size:128;"`                    // 创建者
	Updator     string `json:"updator" gorm:"size:128;"`                    // 更新者
	Model

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

// TableName implement gorm.Tabler interface
func (Config) TableName() string {
	return "sys_config"
}

// ConfigDB config db scope
func ConfigDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Config{})
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

// CConfig 实例
var CConfig = new(cConfig)

// GetWithKey 通过键获取
func (cConfig) GetWithKey(ctx context.Context, key string) (item Config, err error) {
	err = dao.DB.Scopes(ConfigDB(ctx)).
		Where("config_key=?", key).First(&item).Error
	return
}

// 获取 Config
func (cConfig) Get(ctx context.Context, id int) (item Config, err error) {
	err = dao.DB.Scopes(ConfigDB(ctx)).
		Where("config_id=?", id).First(&item).Error
	return
}

// QueryPage 查询
func (cConfig) QueryPage(ctx context.Context, qp ConfigQueryParam) ([]Config, paginator.Info, error) {
	var items []Config

	db := dao.DB.Scopes(ConfigDB(ctx))
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
	db = db.Scopes(DataScope("sys_config", jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// Config 创建
func (cConfig) Create(ctx context.Context, item Config) (Config, error) {
	var i int64

	dao.DB.Scopes(ConfigDB(ctx)).
		Where("config_name=? or config_key=?", item.ConfigName, item.ConfigKey).
		Count(&i)
	if i > 0 {
		return item, iorm.ErrObjectAlreadyExist
	}

	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(ConfigDB(ctx)).Create(&item).Error
	return item, err

}

// Update 更新
func (cConfig) Update(ctx context.Context, id int, up Config) (item Config, err error) {
	if err = dao.DB.Scopes(ConfigDB(ctx)).
		Where("config_id=?", id).First(&item).Error; err != nil {
		return
	}

	if up.ConfigName != "" && up.ConfigName != item.ConfigName {
		return item, errors.New("参数名称不允许修改！")
	}
	if up.ConfigKey != "" && up.ConfigKey != item.ConfigKey {
		return item, errors.New("参数键名不允许修改！")
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	err = dao.DB.Scopes(ConfigDB(ctx)).Model(&item).Updates(&up).Error
	return
}

// Delete 删除
func (cConfig) Delete(ctx context.Context, id int) (err error) {
	return dao.DB.Scopes(ConfigDB(ctx)).
		Where("config_id=?", id).Delete(&Config{}).Error
}

// BatchDelete 批量删除
func (cConfig) BatchDelete(ctx context.Context, id []int) error {
	return dao.DB.Scopes(ConfigDB(ctx)).
		Where("config_id in (?)", id).Delete(&Config{}).Error
}
