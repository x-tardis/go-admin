package models

import (
	"context"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type Content struct {
	Id      int    `json:"id" gorm:"type:int(11);primary_key;auto_increment"` // id
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

func (Content) TableName() string {
	return "sys_content"
}

func ContentDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(Content{})
	}
}

type ContentQueryParam struct {
	CateId string `form:"cateId"`
	Name   string `form:"name"`
	Status string `form:"status"`
	paginator.Param
}

type cContent struct{}

var CContent = new(cContent)

// 获取SysContent带分页
func (cContent) QueryPage(ctx context.Context, qp ContentQueryParam) ([]Content, paginator.Info, error) {
	var err error
	var items []Content

	db := dao.DB.Scopes(ContentDB())
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
	db, err = GetDataScope("sys_content", jwtauth.FromUserId(ctx), db)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// 获取SysContent
func (cContent) Get(_ context.Context, id int) (item Content, err error) {
	err = dao.DB.Scopes(ContentDB()).
		Where("id = ?", id).
		First(&item).Error
	return

}

// 创建SysContent
func (cContent) Create(ctx context.Context, item Content) (Content, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(ContentDB()).Create(&item).Error
	return item, err
}

// 更新SysContent
func (cContent) Update(ctx context.Context, id int, up Content) (item Content, err error) {
	up.Updator = jwtauth.FromUserIdStr(ctx)
	if err = dao.DB.Scopes(ContentDB()).
		Where("id=?", id).First(&item).Error; err != nil {
		return
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = dao.DB.Scopes(ContentDB()).
		Model(&item).Updates(&up).Error
	return
}

// 删除SysContent
func (cContent) Delete(_ context.Context, id int) error {
	return dao.DB.Scopes(ContentDB()).
		Where("id = ?", id).Delete(&Content{}).Error
}

// 批量删除
func (cContent) BatchDelete(_ context.Context, ids []int) error {
	return dao.DB.Scopes(ContentDB()).
		Where("id in (?)", ids).Delete(&Content{}).Error
}
