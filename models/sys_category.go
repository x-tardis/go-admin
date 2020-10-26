package models

import (
	"context"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/trans"
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
	db = db.Scopes(DataScope("sys_category", jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// 获取SysCategory
func (cCategory) Get(ctx context.Context, id int) (item Category, err error) {
	err = dao.DB.Scopes(CategoryDB(ctx)).
		Where("id=?", id).First(&item).Error
	return
}

// 创建SysCategory
func (cCategory) Create(ctx context.Context, item Category) (Category, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(CategoryDB(ctx)).Create(&item).Error
	return item, err
}

// 更新SysCategory
func (cCategory) Update(ctx context.Context, id int, up Category) (update Category, err error) {
	up.Updator = jwtauth.FromUserIdStr(ctx)
	err = dao.DB.Scopes(CategoryDB(ctx)).Where("id=?", id).First(&update).Error
	if err != nil {
		return
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = dao.DB.Scopes(CategoryDB(ctx)).Model(&update).Updates(&up).Error
	return
}

// 删除SysCategory
func (cCategory) Delete(ctx context.Context, id int) error {
	return dao.DB.Scopes(CategoryDB(ctx)).
		Where("id=?", id).Delete(&Category{}).Error
}

// 批量删除
func (cCategory) BatchDelete(ctx context.Context, ids []int) error {
	return dao.DB.Scopes(CategoryDB(ctx)).
		Where("id in (?)", ids).Delete(&Category{}).Error
}
