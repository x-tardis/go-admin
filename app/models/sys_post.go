package models

import (
	"context"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type Post struct {
	PostId   int    `gorm:"primary_key;AUTO_INCREMENT" json:"postId"` // 主键
	PostName string `gorm:"size:128;" json:"postName"`                // 岗位名称
	PostCode string `gorm:"size:128;" json:"postCode"`                // 岗位代码
	Sort     int    `gorm:"" json:"sort"`                             // 岗位排序
	Status   string `gorm:"size:4;" json:"status"`                    // 状态
	Remark   string `gorm:"size:255;" json:"remark"`                  // 备注
	Creator  string `gorm:"size:128;" json:"creator"`
	Updator  string `gorm:"size:128;" json:"updator"`
	Model

	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
}

func (Post) TableName() string {
	return "sys_post"
}

func PostDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(Post{})
	}
}

type PostQueryParam struct {
	PostId   int    `form:"postId"`
	PostName string `form:"postName"`
	PostCode string `form:"postCode"`
	Status   string `form:"status"`
	paginator.Param
}

type cPost struct{}

var CPost = new(cPost)

func (cPost) Query(_ context.Context, qp PostQueryParam) (items []Post, err error) {
	db := deployed.DB.Scopes(PostDB())
	if qp.PostId != 0 {
		db = db.Where("post_id=?", qp.PostId)
	}
	if qp.PostName != "" {
		db = db.Where("post_name=?", qp.PostName)
	}
	if qp.PostCode != "" {
		db = db.Where("post_code=?", qp.PostCode)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}
	err = db.Find(&items).Error
	return
}

func (cPost) QueryPage(ctx context.Context, qp PostQueryParam) ([]Post, paginator.Info, error) {
	var items []Post

	db := deployed.DB.Scopes(PostDB())
	if qp.PostId != 0 {
		db = db.Where("post_id=?", qp.PostId)
	}
	if qp.PostName != "" {
		db = db.Where("post_name=?", qp.PostName)
	}
	if qp.PostCode != "" {
		db = db.Where("post_code=?", qp.PostCode)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId = jwtauth.FromUserId(ctx)
	db, err := dataPermission.GetDataScope("sys_post", db)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

func (cPost) Create(ctx context.Context, item Post) (Post, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := deployed.DB.Scopes(PostDB()).Create(&item).Error
	return item, err
}

func (cPost) Get(_ context.Context, id int) (item Post, err error) {
	err = deployed.DB.Scopes(PostDB()).
		Where("post_id=?", id).First(&item).Error
	return
}

func (cPost) Update(ctx context.Context, id int, up Post) (item Post, err error) {
	up.Updator = jwtauth.FromUserIdStr(ctx)
	if err = deployed.DB.Scopes(PostDB()).First(&item, id).Error; err != nil {
		return
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = deployed.DB.Scopes(PostDB()).Model(&item).Updates(&up).Error
	return
}

func (cPost) Delete(_ context.Context, id int) (err error) {
	return deployed.DB.Scopes(PostDB()).
		Where("post_id=?", id).Delete(&Post{}).Error
}

func (cPost) BatchDelete(_ context.Context, id []int) error {
	return deployed.DB.Scopes(PostDB()).
		Where("post_id in (?)", id).Delete(&Post{}).Error
}
