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

// DictData 字典数据, 字典类型下的实际字典数据
type DictData struct {
	DictId    int    `gorm:"primary_key;auto_increment;" json:"dictId" example:"1"` // 主键
	DictLabel string `gorm:"size:128;" json:"dictLabel"`                            // 标签
	DictValue string `gorm:"size:255;" json:"dictValue"`                            // 值
	DictType  string `gorm:"size:64;" json:"dictType"`                              // 类型
	Sort      int    `gorm:"" json:"sort"`                                          // 排序
	CssClass  string `gorm:"size:128;" json:"cssClass"`                             // (未用)
	ListClass string `gorm:"size:128;" json:"listClass"`                            // (未用)
	IsDefault string `gorm:"size:8;" json:"isDefault"`                              // (未用)
	Default   string `gorm:"size:8;" json:"default"`                                // (未用)
	Status    string `gorm:"size:4;" json:"status"`                                 // 状态
	Remark    string `gorm:"size:255;" json:"remark"`                               // 备注
	Creator   string `gorm:"size:64;" json:"creator"`                               // 创建者
	Updator   string `gorm:"size:64;" json:"updator"`                               // 更新者
	Model

	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
}

// TableName implement schema.Tabler interface
func (DictData) TableName() string {
	return "sys_dict_data"
}

// DictDataDB dict data db
func DictDataDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(DictData{})
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
var CDictData = cDictData{}

// QueryPage 查询,分页
func (cDictData) QueryPage(ctx context.Context, qp DictDataQueryParam) ([]DictData, paginator.Info, error) {
	var items []DictData

	db := dao.DB.Scopes(DictDataDB(ctx))
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
	db = db.Scopes(DataScope(DictData{}, jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, paginator.Info{}, err
	}

	db = db.Order("sort")
	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

func (cDictData) Create(ctx context.Context, item DictData) (DictData, error) {
	var i int64

	if err := dao.DB.Scopes(DictDataDB(ctx)).
		Where("dict_type=?", item.DictType).
		Where("dict_label=? or (dict_label=? and dict_value=?)", item.DictLabel, item.DictLabel, item.DictValue).
		Count(&i).Error; err != nil {
		return item, err
	}
	if i > 0 {
		return item, errors.New("字典标签或者字典键值已经存在！")
	}

	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(DictDataDB(ctx)).Create(&item).Error
	return item, err
}

// Get 通过dictCode(主键)
func (cDictData) Get(ctx context.Context, dictId int) (item DictData, err error) {
	err = dao.DB.Scopes(DictDataDB(ctx)).
		Where("dict_id=?", dictId).First(&item).Error
	return
}

// GetWithType 通过dictType获取
func (cDictData) GetWithType(ctx context.Context, dictType string) (items []DictData, err error) {
	err = dao.DB.Scopes(DictDataDB(ctx)).
		Where("dict_type = ?", dictType).
		Order("sort").Find(&items).Error
	return
}

// Update 更新
func (cDictData) Update(ctx context.Context, id int, up DictData) (item DictData, err error) {
	if err = dao.DB.Scopes(DictDataDB(ctx)).
		Where("dict_id = ?", id).First(&item).Error; err != nil {
		return
	}

	if up.DictLabel != "" && up.DictLabel != item.DictLabel {
		return item, errors.New("标签不允许修改！")
	}

	if up.DictValue != "" && up.DictValue != item.DictValue {
		return item, errors.New("键值不允许修改！")
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	err = dao.DB.Scopes(DictDataDB(ctx)).Model(&item).Updates(&up).Error
	return
}

// Delete 删除
func (cDictData) Delete(ctx context.Context, id int) error {
	return dao.DB.Scopes(DictDataDB(ctx)).
		Where("dict_id = ?", id).Delete(&DictData{}).Error
}

// BatchDelete 批量删除
func (cDictData) BatchDelete(ctx context.Context, id []int) error {
	return dao.DB.Scopes(DictDataDB(ctx)).
		Where("dict_id in (?)", id).Delete(&DictData{}).Error
}
