package models

import (
	"context"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type Category struct {
	Id      int    `json:"id" gorm:"type:int(11);primary_key;AUTO_INCREMENT"` // 主键
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

func (Category) TableName() string {
	return "sys_category"
}

func CategoryDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Category{})
	}
}

type CategoryQueryParam struct {
	Name   string `form:"name"`
	Status string `form:"status"`
	paginator.Param
}

type cCategory struct{}

var CCategory = cCategory{}

// 获取SysCategory带分页
func (cCategory) QueryPage(ctx context.Context, qp CategoryQueryParam) ([]Category, paginator.Info, error) {
	var items []Category

	db := dao.DB.Scopes(CategoryDB(ctx))
	if qp.Name != "" {
		db = db.Where("name=?", qp.Name)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	db = db.Scopes(DataScope(Category{}, jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// 获取SysCategory
func (cCategory) Get(ctx context.Context, id int) (item Category, err error) {
	err = dao.DB.Scopes(CategoryDB(ctx)).
		First(&item, "id=?", id).Error
	return
}

// 创建SysCategory
func (cCategory) Create(ctx context.Context, item Category) (Category, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(CategoryDB(ctx)).Create(&item).Error
	return item, err
}

// 更新SysCategory
func (sf cCategory) Update(ctx context.Context, id int, up Category) error {
	oldItem, err := sf.Get(ctx, id)
	if err != nil {
		return err
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	return dao.DB.Scopes(CategoryDB(ctx)).Model(&oldItem).Updates(&up).Error
}

// 删除SysCategory
func (cCategory) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return nil
	}
	return dao.DB.Scopes(CategoryDB(ctx)).
		Delete(&Category{}, "id=?", id).Error
}

// 批量删除
func (sf cCategory) BatchDelete(ctx context.Context, ids []int) error {
	switch len(ids) {
	case 0:
		return nil
	case 1:
		return sf.Delete(ctx, ids[0])
	default:
		return dao.DB.Scopes(CategoryDB(ctx)).
			Delete(&Category{}, "id in (?)", ids).Error
	}
}
