package models

import (
	"context"
	"errors"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/app/models/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// DictType 字典类型
type DictType struct {
	DictId   int    `gorm:"primary_key;auto_increment;" json:"dictId"` // 主键
	DictName string `gorm:"size:128;" json:"dictName"`                 // 名称
	DictType string `gorm:"size:128;" json:"dictType"`                 // 类型
	Status   string `gorm:"size:4;" json:"status"`                     // 状态
	Remark   string `gorm:"size:255;" json:"remark"`                   // 备注
	Creator  string `gorm:"size:11;" json:"creator"`                   // 创建者
	Updator  string `gorm:"size:11;" json:"updator"`                   // 更新者
	Model

	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
}

// TableName implement gorm.Tabler interface
func (DictType) TableName() string {
	return "sys_dict_type"
}

// DictTypeDB dict type db
func DictTypeDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(DictType{})
	}
}

// DictTypeQueryParam 查询参数
type DictTypeQueryParam struct {
	DictId   int    `form:"dictId"`
	DictName string `form:"dictName"`
	DictType string `form:"dictType"`
	paginator.Param
}

type cDictType struct{}

// CDictType 实例
var CDictType = new(cDictType)

// Query 查询,非分页
func (cDictType) Query(_ context.Context, qp DictTypeQueryParam) ([]DictType, error) {
	var item []DictType

	db := dao.DB.Scopes(DictTypeDB())
	if qp.DictId != 0 {
		db = db.Where("dict_id=?", qp.DictId)
	}
	if qp.DictName != "" {
		db = db.Where("dict_name=?", qp.DictName)
	}
	if qp.DictType != "" {
		db = db.Where("dict_type=?", qp.DictType)
	}

	err := db.Find(&item).Error
	return item, err
}

// QueryPage 查询,分页
func (cDictType) QueryPage(ctx context.Context, qp DictTypeQueryParam) ([]DictType, paginator.Info, error) {
	var err error
	var items []DictType

	db := dao.DB.Scopes(DictTypeDB())
	if qp.DictId != 0 {
		db = db.Where("dict_id=?", qp.DictId)
	}
	if qp.DictName != "" {
		db = db.Where("dict_name=?", qp.DictName)
	}
	if qp.DictType != "" {
		db = db.Where("dict_type=?", qp.DictType)
	}

	// 数据权限控制
	db, err = GetDataScope("sys_dict_type", jwtauth.FromUserId(ctx), db)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// Get 通过id或name查询
func (cDictType) Get(_ context.Context, dictId int, dictName string) (item DictType, err error) {
	db := dao.DB.Scopes(DictTypeDB())
	if dictId != 0 {
		db = db.Where("dict_id=?", dictId)
	}
	if dictName != "" {
		db = db.Where("dict_name=?", dictName)
	}
	err = db.First(&item).Error
	return
}

// Create 创建
func (cDictType) Create(ctx context.Context, item DictType) (DictType, error) {
	var i int64

	dao.DB.Scopes(DictTypeDB()).
		Where("dict_name=? or dict_type=?", item.DictName, item.DictType).
		Count(&i)
	if i > 0 {
		return item, errors.New("字典名称或者字典类型已经存在！")
	}

	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(DictTypeDB()).Create(&item).Error
	return item, err
}

// Update 更新
func (cDictType) Update(ctx context.Context, id int, up DictType) (item DictType, err error) {
	if err = dao.DB.Scopes(DictTypeDB()).First(&item, id).Error; err != nil {
		return
	}

	if up.DictName != "" && up.DictName != item.DictName {
		return item, errors.New("名称不允许修改！")
	}
	if up.DictType != "" && up.DictType != item.DictType {
		return item, errors.New("类型不允许修改！")
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	err = dao.DB.Scopes(DictTypeDB()).Model(&item).Updates(&up).Error
	return
}

// Delete 根据id删除
func (cDictType) Delete(_ context.Context, id int) error {
	return dao.DB.Scopes(DictTypeDB()).
		Where("dict_id=?", id).Delete(&DictData{}).Error
}

// BatchDelete 根据id列表批量删除
func (cDictType) BatchDelete(_ context.Context, ids []int) error {
	return dao.DB.Scopes(DictTypeDB()).
		Where("dict_id in (?)", ids).Delete(&DictType{}).Error
}
