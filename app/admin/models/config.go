package models

import (
	"errors"
	"strconv"
	_ "time"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

type SysConfig struct {
	ConfigId    int    `json:"configId" gorm:"primary_key;auto_increment;"` //编码
	ConfigName  string `json:"configName" gorm:"size:128;"`                 //参数名称
	ConfigKey   string `json:"configKey" gorm:"size:128;"`                  //参数键名
	ConfigValue string `json:"configValue" gorm:"size:255;"`                //参数键值
	ConfigType  string `json:"configType" gorm:"size:64;"`                  //是否系统内置
	Remark      string `json:"remark" gorm:"size:128;"`                     //备注
	CreateBy    string `json:"createBy" gorm:"size:128;"`
	UpdateBy    string `json:"updateBy" gorm:"size:128;"`
	BaseModel

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}

// Config 创建
func (e *SysConfig) Create() (SysConfig, error) {
	var doc SysConfig
	var i int64
	deployed.DB.Table(e.TableName()).Where("config_name=? or config_key = ?", e.ConfigName, e.ConfigKey).Count(&i)
	if i > 0 {
		return doc, errors.New("参数名称或者参数键名已经存在！")
	}

	result := deployed.DB.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取 Config
func (e *SysConfig) Get() (SysConfig, error) {
	var doc SysConfig

	table := deployed.DB.Table(e.TableName())
	if e.ConfigId != 0 {
		table = table.Where("config_id = ?", e.ConfigId)
	}

	if e.ConfigKey != "" {
		table = table.Where("config_key = ?", e.ConfigKey)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *SysConfig) GetPage(pageSize int, pageIndex int) ([]SysConfig, int, error) {
	var doc []SysConfig

	table := deployed.DB.Table(e.TableName())

	if e.ConfigName != "" {
		table = table.Where("config_name = ?", e.ConfigName)
	}
	if e.ConfigKey != "" {
		table = table.Where("config_key = ?", e.ConfigKey)
	}
	if e.ConfigType != "" {
		table = table.Where("config_type = ?", e.ConfigType)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = strconv.Atoi(e.DataScope)
	table, err := dataPermission.GetDataScope("sys_config", table)
	if err != nil {
		return nil, 0, err
	}
	var count int64

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	//table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, int(count), nil
}

func (e *SysConfig) Update(id int) (update SysConfig, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("config_id = ?", id).First(&update).Error; err != nil {
		return
	}

	if e.ConfigName != "" && e.ConfigName != update.ConfigName {
		return update, errors.New("参数名称不允许修改！")
	}

	if e.ConfigKey != "" && e.ConfigKey != update.ConfigKey {
		return update, errors.New("参数键名不允许修改！")
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = deployed.DB.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *SysConfig) Delete() (success bool, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("config_id = ?", e.ConfigId).Delete(&SysConfig{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *SysConfig) BatchDelete(id []int) (Result bool, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("config_id in (?)", id).Delete(&SysConfig{}).Error; err != nil {
		return
	}
	Result = true
	return
}
