package models

import (
	"context"
	"errors"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// Config 参数配置
// 主要用于系统参数配置,包含名称,键名,键值,是否系统内置
// 可用于web一些系统化配置,如框架样式,主框架侧边主题,用户管理默认密码等
// 由开发者动态决定,以适配不同的需求.
type Config struct {
	ConfigId    int    `json:"configId" gorm:"primary_key;auto_increment;"` // 主键
	ConfigName  string `json:"configName" gorm:"size:128;"`                 // 名称
	ConfigKey   string `json:"configKey" gorm:"size:128;"`                  // 键名
	ConfigValue string `json:"configValue" gorm:"size:255;"`                // 键值
	ConfigType  string `json:"configType" gorm:"size:64;"`                  // 是否系统内置(Y/N) 字典sys_yes_no
	Remark      string `json:"remark" gorm:"size:128;"`                     // 备注
	Creator     string `json:"creator" gorm:"size:128;"`                    // 创建者
	Updator     string `json:"updator" gorm:"size:128;"`                    // 更新者
	Model

	DataScope string `json:"-" gorm:"-"`
	Params    string `json:"-"  gorm:"-"`
}

// TableName implement schema.Tabler interface
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
var CConfig = cConfig{}

// GetWithKey 通过键获取
func (cConfig) GetWithKey(ctx context.Context, key string) (item Config, err error) {
	err = dao.DB.Scopes(ConfigDB(ctx)).
		First(&item, "config_key=?", key).Error
	return
}

// 获取 Config
func (cConfig) Get(ctx context.Context, id int) (item Config, err error) {
	err = dao.DB.Scopes(ConfigDB(ctx)).
		First(&item, "config_id=?", id).Error
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
	db = db.Scopes(DataScope(Config{}, jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// Config 创建
func (cConfig) Create(ctx context.Context, item Config) (Config, error) {
	var count int64

	// 键值应
	dao.DB.Scopes(ConfigDB(ctx)).
		Where("config_name=?", item.ConfigName).
		Or("config_key=?", item.ConfigKey).
		Count(&count)
	if count > 0 {
		return item, iorm.ErrObjectAlreadyExist
	}

	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(ConfigDB(ctx)).Create(&item).Error
	return item, err

}

// Update 更新
func (sf cConfig) Update(ctx context.Context, id int, up Config) error {
	item, err := sf.Get(ctx, id)
	if err != nil {
		return err
	}

	if up.ConfigName != "" && up.ConfigName != item.ConfigName {
		return errors.New("参数名称不允许修改！")
	}
	if up.ConfigKey != "" && up.ConfigKey != item.ConfigKey {
		return errors.New("参数键名不允许修改！")
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	return dao.DB.Scopes(ConfigDB(ctx)).Model(&item).Updates(&up).Error
}

// Delete 删除
func (cConfig) Delete(ctx context.Context, id int) (err error) {
	return dao.DB.Scopes(ConfigDB(ctx)).
		Delete(&Config{}, "config_id=?", id).Error
}

// BatchDelete 批量删除
func (sf cConfig) BatchDelete(ctx context.Context, ids []int) error {
	switch len(ids) {
	case 0:
		return nil
	case 1:
		return sf.Delete(ctx, ids[0])
	default:
		return dao.DB.Scopes(ConfigDB(ctx)).
			Delete(&Config{}, "config_id in (?)", ids).Error
	}
}
