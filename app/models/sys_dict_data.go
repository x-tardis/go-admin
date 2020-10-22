package models

import (
	"context"
	"errors"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// DictData 字典数据
type DictData struct {
	DictCode  int    `gorm:"primary_key;auto_increment;" json:"dictCode" example:"1"` // 主键
	DictLabel string `gorm:"size:128;" json:"dictLabel"`                              // 标签
	DictValue string `gorm:"size:255;" json:"dictValue"`                              // 值
	DictType  string `gorm:"size:64;" json:"dictType"`                                // 类型
	DictSort  int    `gorm:"" json:"dictSort"`                                        // 排序
	CssClass  string `gorm:"size:128;" json:"cssClass"`                               // (未用)
	ListClass string `gorm:"size:128;" json:"listClass"`                              // (未用)
	IsDefault string `gorm:"size:8;" json:"isDefault"`                                // (未用)
	Status    string `gorm:"size:4;" json:"status"`                                   // 状态
	Default   string `gorm:"size:8;" json:"default"`                                  // (未用)
	Remark    string `gorm:"size:255;" json:"remark"`                                 // 备注
	Creator   string `gorm:"size:64;" json:"creator"`                                 // 创建者
	Updator   string `gorm:"size:64;" json:"updator"`                                 // 更新者
	Model

	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
}

// TableName implement gorm.Tabler interface
func (DictData) TableName() string {
	return "sys_dict_data"
}

// DictDataDB dict data db
func DictDataDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(DictData{})
	}
}

// DictDataQueryParam 查询参数
type DictDataQueryParam struct {
	DictLabel string `form:"dictLabel"`
	DictType  string `form:"dictType"`
	Status    string `form:"status"`
	paginator.Param
}

type cDictData struct{}

// CDictData 实例
var CDictData = new(cDictData)

// QueryPage 查询,分页
func (cDictData) QueryPage(ctx context.Context, qp DictDataQueryParam) ([]DictData, paginator.Info, error) {
	var err error
	var items []DictData

	db := deployed.DB.Scopes(DictDataDB())
	if qp.DictType != "" {
		db = db.Where("dict_type=?", qp.DictType)
	}
	if qp.DictLabel != "" {
		db = db.Where("dict_label=?", qp.DictLabel)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}

	// 数据权限控制
	db, err = GetDataScope("sys_dict_data", jwtauth.FromUserId(ctx), db)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

func (cDictData) Create(ctx context.Context, item DictData) (DictData, error) {
	var i int64

	if err := deployed.DB.Scopes(DictDataDB()).
		Where("dict_type=?", item.DictType).
		Where("dict_label=? or (dict_label=? and dict_value=?)", item.DictLabel, item.DictLabel, item.DictValue).
		Count(&i).Error; err != nil {
		return item, err
	}
	if i > 0 {
		return item, errors.New("字典标签或者字典键值已经存在！")
	}

	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := deployed.DB.Scopes(DictDataDB()).Create(&item).Error
	return item, err
}

// Get 通过dictCode(主键)
func (cDictData) Get(_ context.Context, dictCode int) (item DictData, err error) {
	err = deployed.DB.Scopes(DictDataDB()).
		Where("dict_code=?", dictCode).First(&item).Error
	return
}

// GetWithType 通过dictType获取
func (cDictData) GetWithType(_ context.Context, dictType string) (items []DictData, err error) {
	err = deployed.DB.Scopes(DictDataDB()).
		Where("dict_type = ?", dictType).
		Order("dict_sort").Find(&items).Error
	return
}

// Update 更新
func (cDictData) Update(ctx context.Context, id int, up DictData) (item DictData, err error) {
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

	up.Updator = jwtauth.FromUserIdStr(ctx)
	err = deployed.DB.Scopes(DictDataDB()).Model(&item).Updates(&up).Error
	return
}

// Delete 删除
func (cDictData) Delete(_ context.Context, id int) error {
	return deployed.DB.Scopes(DictDataDB()).
		Where("dict_code = ?", id).Delete(&DictData{}).Error
}

// BatchDelete 批量删除
func (cDictData) BatchDelete(_ context.Context, id []int) error {
	return deployed.DB.Scopes(DictDataDB()).
		Where("dict_code in (?)", id).Delete(&DictData{}).Error
}
