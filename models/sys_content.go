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

// Content 内容
type Content struct {
	Id      int    `json:"id" gorm:"type:int(11);primary_key;auto_increment"` // 主键
	CateId  string `json:"cateId" gorm:"type:int(11);"`                       // 分类id
	Name    string `json:"name" gorm:"type:varchar(255);"`                    // 名称
	Status  string `json:"status" gorm:"type:int(1);"`                        // 状态
	Img     string `json:"img" gorm:"type:varchar(255);"`                     // 图片
	Content string `json:"content" gorm:"type:text;"`                         // 内容
	Remark  string `json:"remark" gorm:"type:varchar(255);"`                  // 备注
	Sort    string `json:"sort" gorm:"type:int(4);"`                          // 排序
	Creator string `json:"creator" gorm:"type:varchar(128);"`                 // 创建者
	Updator string `json:"updator" gorm:"type:varchar(128);"`                 // 更新者
	Model

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

// TableName implement schema.Tabler interface
func (Content) TableName() string {
	return "sys_content"
}

// ContentDB content db scopes
func ContentDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Content{})
	}
}

// ContentQueryParam 查询参数
type ContentQueryParam struct {
	CateId string `form:"cateId"`
	Name   string `form:"name"`
	Status string `form:"status"`
	paginator.Param
}

type cContent struct{}

// CContent 实例
var CContent = cContent{}

// QueryPage 查询,带分页
func (cContent) QueryPage(ctx context.Context, qp ContentQueryParam) ([]Content, paginator.Info, error) {
	var err error
	var items []Content

	db := dao.DB.Scopes(ContentDB(ctx))
	if qp.CateId != "" {
		db = db.Where("cate_id = ?", qp.CateId)
	}
	if qp.Name != "" {
		db = db.Where("name like ?", "%"+qp.Name+"%")
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	db = db.Scopes(DataScope(Content{}, jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// Get 获取
func (cContent) Get(ctx context.Context, id int) (item Content, err error) {
	err = dao.DB.Scopes(ContentDB(ctx)).
		Where("id=?", id).
		First(&item).Error
	return

}

// Create 创建
func (cContent) Create(ctx context.Context, item Content) (Content, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(ContentDB(ctx)).Create(&item).Error
	return item, err
}

// 更新SysContent
func (cContent) Update(ctx context.Context, id int, up Content) (item Content, err error) {
	if err = dao.DB.Scopes(ContentDB(ctx)).
		Where("id=?", id).First(&item).Error; err != nil {
		return
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	err = dao.DB.Scopes(ContentDB(ctx)).
		Model(&item).Updates(&up).Error
	return
}

// Delete 删除
func (cContent) Delete(ctx context.Context, id int) error {
	return dao.DB.Scopes(ContentDB(ctx)).
		Where("id = ?", id).Delete(&Content{}).Error
}

// BatchDelete 批量删除
func (cContent) BatchDelete(ctx context.Context, ids []int) error {
	return dao.DB.Scopes(ContentDB(ctx)).
		Where("id in (?)", ids).Delete(&Content{}).Error
}
