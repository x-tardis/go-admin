package models

import (
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

type SysCategory struct {
	Id      int    `json:"id" gorm:"type:int(11);primary_key;AUTO_INCREMENT"` // 分类Id
	Name    string `json:"name" gorm:"type:varchar(255);"`                    // 名称
	Img     string `json:"img" gorm:"type:varchar(255);"`                     // 图片
	Sort    string `json:"sort" gorm:"type:int(4);"`                          // 排序
	Status  string `json:"status" gorm:"type:int(1);"`                        // 状态
	Remark  string `json:"remark" gorm:"type:varchar(255);"`                  // 备注
	Creator string `json:"creator" gorm:"type:varchar(64);"`                  // 创建者
	Updator string `json:"updator" gorm:"type:varchar(64);"`                  // 更新者
	Model

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (SysCategory) TableName() string {
	return "sys_category"
}

// 创建SysCategory
func (e *SysCategory) Create() (SysCategory, error) {
	var doc SysCategory
	result := deployed.DB.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取SysCategory
func (e *SysCategory) Get() (SysCategory, error) {
	var doc SysCategory
	table := deployed.DB.Table(e.TableName())

	if e.Id != 0 {
		table = table.Where("id = ?", e.Id)
	}

	if e.Name != "" {
		table = table.Where("name = ?", e.Name)
	}

	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取SysCategory带分页
func (e *SysCategory) GetPage(param paginator.Param) ([]SysCategory, paginator.Info, error) {
	var doc []SysCategory

	table := deployed.DB.Table(e.TableName())

	if e.Name != "" {
		table = table.Where("name = ?", e.Name)
	}

	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(DataPermission)
	dataPermission.UserId = cast.ToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, paginator.Info{}, err
	}
	ifc, err := iorm.QueryPages(table, param, &doc)
	if err != nil {
		return nil, ifc, nil
	}
	return doc, ifc, nil
}

// 更新SysCategory
func (e *SysCategory) Update(id int) (update SysCategory, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = deployed.DB.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除SysCategory
func (e *SysCategory) Delete(id int) (success bool, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("id = ?", id).Delete(&SysCategory{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *SysCategory) BatchDelete(id []int) (Result bool, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("id in (?)", id).Delete(&SysCategory{}).Error; err != nil {
		return
	}
	Result = true
	return
}
