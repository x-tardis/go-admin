package models

import (
	"context"
	"errors"

	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type DictData struct {
	DictCode  int    `gorm:"primary_key;auto_increment;" json:"dictCode" example:"1"` // 字典编码
	DictSort  int    `gorm:"" json:"dictSort"`                                        // 显示顺序
	DictLabel string `gorm:"size:128;" json:"dictLabel"`                              // 数据标签
	DictValue string `gorm:"size:255;" json:"dictValue"`                              // 数据键值
	DictType  string `gorm:"size:64;" json:"dictType"`                                // 字典类型
	CssClass  string `gorm:"size:128;" json:"cssClass"`                               //
	ListClass string `gorm:"size:128;" json:"listClass"`                              //
	IsDefault string `gorm:"size:8;" json:"isDefault"`                                //
	Status    string `gorm:"size:4;" json:"status"`                                   // 状态
	Default   string `gorm:"size:8;" json:"default"`                                  //
	Remark    string `gorm:"size:255;" json:"remark"`                                 // 备注
	CreateBy  string `gorm:"size:64;" json:"createBy"`                                //
	UpdateBy  string `gorm:"size:64;" json:"updateBy"`                                //
	Model

	Params    string `gorm:"-" json:"params"`
	DataScope string `gorm:"-" json:"dataScope"`
}

func (DictData) TableName() string {
	return "sys_dict_data"
}

func DictDataDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(DictData{})
	}
}

type CallDictData struct{}

func (e *DictData) Create() (DictData, error) {
	var doc DictData

	var i int64
	if err := deployed.DB.Table(e.TableName()).
		Where("dict_type = ?", e.DictType).
		Where("dict_label=? or (dict_label=? and dict_value = ?)", e.DictLabel, e.DictLabel, e.DictValue).
		Count(&i).Error; err != nil {
		return doc, err
	}
	if i > 0 {
		return doc, errors.New("字典标签或者字典键值已经存在！")
	}

	result := deployed.DB.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (CallDictData) Get(_ context.Context, dictCode int, dictLabel string) (item DictData, err error) {
	db := deployed.DB.Scopes(DictDataDB())
	if dictCode != 0 {
		db = db.Where("dict_code = ?", dictCode)
	}
	if dictLabel != "" {
		db = db.Where("dict_label = ?", dictLabel)
	}
	err = db.First(&item).Error
	return
}

func (CallDictData) GetWithType(_ context.Context, dictType string) (items []DictData, err error) {
	err = deployed.DB.Scopes(DictDataDB()).
		Where("dict_type = ?", dictType).
		Order("dict_sort").Find(&items).Error
	return
}

func (e *DictData) GetPage(param paginator.Param) ([]DictData, paginator.Info, error) {
	var doc []DictData

	table := deployed.DB.Table(e.TableName())
	if e.DictCode != 0 {
		table = table.Where("dict_code = ?", e.DictCode)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}
	if e.DictLabel != "" {
		table = table.Where("dict_label = ?", e.DictLabel)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId = cast.ToInt(e.DataScope)
	table, err := dataPermission.GetDataScope("sys_dict_data", table)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	ifc, err := iorm.QueryPages(table, param, doc)
	if err != nil {
		return nil, ifc, err
	}
	return doc, ifc, nil
}

func (CallDictData) Update(ctx context.Context, id int, up DictData) (item DictData, err error) {
	up.UpdateBy = jwtauth.FromUserIdStr(ctx)
	if err = deployed.DB.Scopes(DictDataDB()).
		Where("dict_code = ?", id).First(&item).Error; err != nil {
		return
	}

	if up.DictLabel != "" && up.DictLabel != item.DictLabel {
		return item, errors.New("标签不允许修改！")
	}

	if up.DictValue != "" && up.DictValue != item.DictValue {
		return item, errors.New("键值不允许修改！")
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = deployed.DB.Scopes(DictDataDB()).Model(&item).Updates(&up).Error
	return
}

func (CallDictData) Delete(_ context.Context, id int) error {
	return deployed.DB.Scopes(DictDataDB()).
		Where("dict_code = ?", id).Delete(&DictData{}).Error
}

func (CallDictData) BatchDelete(_ context.Context, id []int) error {
	return deployed.DB.Scopes(DictDataDB()).
		Where("dict_code in (?)", id).Delete(&DictData{}).Error
}
