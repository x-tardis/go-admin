package models

import (
	"errors"

	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

type DictType struct {
	DictId   int    `gorm:"primary_key;auto_increment;" json:"dictId"`
	DictName string `gorm:"size:128;" json:"dictName"` //字典名称
	DictType string `gorm:"size:128;" json:"dictType"` //字典类型
	Status   string `gorm:"size:4;" json:"status"`     //状态
	CreateBy string `gorm:"size:11;" json:"createBy"`  //创建者
	UpdateBy string `gorm:"size:11;" json:"updateBy"`  //更新者
	Remark   string `gorm:"size:255;" json:"remark"`   //备注
	Model

	DataScope string `gorm:"-" json:"dataScope"` //
	Params    string `gorm:"-" json:"params"`    //
}

func (DictType) TableName() string {
	return "sys_dict_type"
}

func (e *DictType) Create() (DictType, error) {
	var doc DictType

	var i int64
	deployed.DB.Table(e.TableName()).Where("dict_name=? or dict_type = ?", e.DictName, e.DictType).Count(&i)
	if i > 0 {
		return doc, errors.New("字典名称或者字典类型已经存在！")
	}

	result := deployed.DB.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *DictType) Get() (DictType, error) {
	var doc DictType

	table := deployed.DB.Table(e.TableName())
	if e.DictId != 0 {
		table = table.Where("dict_id = ?", e.DictId)
	}
	if e.DictName != "" {
		table = table.Where("dict_name = ?", e.DictName)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *DictType) GetList() ([]DictType, error) {
	var doc []DictType

	table := deployed.DB.Table(e.TableName())
	if e.DictId != 0 {
		table = table.Where("dict_id = ?", e.DictId)
	}
	if e.DictName != "" {
		table = table.Where("dict_name = ?", e.DictName)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}

	if err := table.Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *DictType) GetPage(param paginator.Param) ([]DictType, paginator.Info, error) {
	var doc []DictType
	var db *gorm.DB

	table := deployed.DB.Table(e.TableName())
	if e.DictId != 0 {
		db = db.Where("dict_id = ?", e.DictId)
	}
	if e.DictName != "" {
		db = db.Where("dict_name = ?", e.DictName)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId = cast.ToInt(e.DataScope)
	table, err := dataPermission.GetDataScope("sys_dict_type", table)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	ifc, err := iorm.QueryPages(table, param, &doc)
	if err != nil {
		return nil, ifc, nil
	}
	return doc, ifc, nil
}

func (e *DictType) Update(id int) (update DictType, err error) {
	if err = deployed.DB.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	if e.DictName != "" && e.DictName != update.DictName {
		return update, errors.New("名称不允许修改！")
	}

	if e.DictType != "" && e.DictType != update.DictType {
		return update, errors.New("类型不允许修改！")
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = deployed.DB.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *DictType) Delete(id int) (success bool, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("dict_id = ?", id).Delete(&DictData{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *DictType) BatchDelete(id []int) (Result bool, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("dict_id in (?)", id).Delete(&DictType{}).Error; err != nil {
		return
	}
	Result = true
	return
}
