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
	PostId   int    `gorm:"primary_key;AUTO_INCREMENT" json:"postId"` // 岗位编号
	PostName string `gorm:"size:128;" json:"postName"`                // 岗位名称
	PostCode string `gorm:"size:128;" json:"postCode"`                // 岗位代码
	Sort     int    `gorm:"" json:"sort"`                             // 岗位排序
	Status   string `gorm:"size:4;" json:"status"`                    // 状态
	Remark   string `gorm:"size:255;" json:"remark"`                  // 描述
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

type CallPost struct{}

func (CallPost) Create(ctx context.Context, item Post) (Post, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := deployed.DB.Scopes(PostDB()).Create(&item).Error
	return item, err
}

func (CallPost) Get(_ context.Context, id int) (item Post, err error) {
	err = deployed.DB.Scopes(PostDB()).
		Where("post_id = ?", id).First(&item).Error
	return
}

func (e *Post) GetList() ([]Post, error) {
	var doc []Post

	table := deployed.DB.Table(e.TableName())
	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}
	if e.PostName != "" {
		table = table.Where("post_name = ?", e.PostName)
	}
	if e.PostCode != "" {
		table = table.Where("post_code = ?", e.PostCode)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (CallPost) QueryPage(ctx context.Context, qp PostQueryParam) ([]Post, paginator.Info, error) {
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

	dataScope := jwtauth.FromUserId(ctx)
	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId = dataScope
	db, err := dataPermission.GetDataScope("sys_post", db)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	ifc, err := iorm.QueryPages(db, qp.Param, &items)
	return items, ifc, err
}

func (CallPost) Update(ctx context.Context, id int, up Post) (item Post, err error) {
	up.Updator = jwtauth.FromUserIdStr(ctx)
	if err = deployed.DB.Scopes(PostDB()).First(&item, id).Error; err != nil {
		return
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = deployed.DB.Scopes(PostDB()).Model(&item).Updates(&up).Error
	return
}

func (CallPost) Delete(_ context.Context, id int) (err error) {
	return deployed.DB.Scopes(PostDB()).
		Where("post_id=?", id).Delete(&Post{}).Error
}

func (CallPost) BatchDelete(_ context.Context, id []int) error {
	return deployed.DB.Scopes(PostDB()).
		Where("post_id in (?)", id).Delete(&Post{}).Error
}
