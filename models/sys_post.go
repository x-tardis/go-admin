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

// Post 岗位
type Post struct {
	PostId   int    `gorm:"primary_key;AUTO_INCREMENT" json:"postId"` // 主键
	PostCode string `gorm:"size:128;" json:"postCode"`                // 编码,比如CEO,CTO
	PostName string `gorm:"size:128;" json:"postName"`                // 名称
	Sort     int    `gorm:"" json:"sort"`                             // 排序
	Status   string `gorm:"size:4;" json:"status"`                    // 状态
	Remark   string `gorm:"size:255;" json:"remark"`                  // 备注
	Creator  string `gorm:"size:128;" json:"creator"`                 // 创建者
	Updator  string `gorm:"size:128;" json:"updator"`                 // 更新者
	Model

	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
}

// TableName implement schema.Tabler interface
func (Post) TableName() string {
	return "sys_post"
}

// PostDB post db scope
func PostDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Post{})
	}
}

// PostQueryParam 查询参数
type PostQueryParam struct {
	PostName string `form:"postName"`
	PostCode string `form:"postCode"`
	Status   string `form:"status"`
	paginator.Param
}

type cPost struct{}

// CPost post 实例
var CPost = cPost{}

// Query 查询岗位信息, 非分页查询
func (cPost) Query(ctx context.Context, qp PostQueryParam) (items []Post, err error) {
	db := dao.DB.Scopes(PostDB(ctx))
	if qp.PostCode != "" {
		db = db.Where("post_code=?", qp.PostCode)
	}
	if qp.PostName != "" {
		db = db.Where("post_name=?", qp.PostName)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}

	err = db.Order("sort").Find(&items).Error
	return
}

// QueryPage 查询岗位信息,分页查询
func (cPost) QueryPage(ctx context.Context, qp PostQueryParam) ([]Post, paginator.Info, error) {
	var items []Post

	db := dao.DB.Scopes(PostDB(ctx))
	if qp.PostCode != "" {
		db = db.Where("post_code=?", qp.PostCode)
	}
	if qp.PostName != "" {
		db = db.Where("post_name=?", qp.PostName)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}
	db = db.Order("sort")

	// 数据权限控制
	db = db.Scopes(DataScope(Post{}, jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// Get 获取岗位信息
func (cPost) Get(ctx context.Context, id int) (item Post, err error) {
	err = dao.DB.Scopes(PostDB(ctx)).
		First(&item, "post_id=?", id).Error
	return
}

// Create 创建岗位信息
func (cPost) Create(ctx context.Context, item Post) (Post, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(PostDB(ctx)).Create(&item).Error
	return item, err
}

// Update 更新岗位信息
func (sf cPost) Update(ctx context.Context, id int, up Post) error {
	oldItem, err := sf.Get(ctx, id)
	if err != nil {
		return err
	}
	up.Updator = jwtauth.FromUserIdStr(ctx)
	return dao.DB.Scopes(PostDB(ctx)).Model(&oldItem).Updates(&up).Error
}

// Delete 删除指定id
func (cPost) Delete(ctx context.Context, id int) (err error) {
	count, _ := CUser.GetCountWithPostId(ctx, id)
	if count > 0 {
		return errors.New("当前岗位存在用户,不能删除")
	}
	return dao.DB.Scopes(PostDB(ctx)).
		Delete(&Post{}, "post_id=?", id).Error
}

// BatchDelete 删除批量id
func (cPost) BatchDelete(ctx context.Context, ids []int) error {
	count, _ := CUser.GetCountWithPostIds(ctx, ids)
	if count > 0 {
		return errors.New("岗位存在用户,不能删除")
	}
	return dao.DB.Scopes(PostDB(ctx)).
		Delete(&Post{}, "post_id in (?)", ids).Error
}
