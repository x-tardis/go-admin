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

type FileInfo struct {
	Id      int    `json:"id"`                                // 主键
	Type    string `json:"type" gorm:"type:varchar(255);"`    // 类型
	Name    string `json:"name" gorm:"type:varchar(255);"`    // 名称
	Size    string `json:"size" gorm:"type:int(11);"`         // 大小
	Pid     int    `json:"pid" gorm:"type:int(11);"`          // 父级id
	Source  string `json:"source" gorm:"type:varchar(255);"`  // 源
	Url     string `json:"url" gorm:"type:varchar(255);"`     // 路径
	FullUrl string `json:"fullUrl" gorm:"type:varchar(255);"` // 全路径
	Creator string `json:"creator" gorm:"type:varchar(128);"` // 创建者
	Updator string `json:"updator" gorm:"type:varchar(128);"` // 更新者
	Model

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

// TableName implement schema.Tabler interface
func (FileInfo) TableName() string {
	return "sys_file_info"
}

// FileInfoDB file info db scopes
func FileInfoDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(FileInfo{})
	}
}

// FileInfoQueryParam 查询参数
type FileInfoQueryParam struct {
	Pid int `form:"pid"`
	paginator.Param
}

type cFileInfo struct{}

// CFileInfo 实例
var CFileInfo = cFileInfo{}

// QueryPage 查询,带分页
func (cFileInfo) QueryPage(ctx context.Context, qp FileInfoQueryParam) ([]FileInfo, paginator.Info, error) {
	var err error
	var items []FileInfo

	db := dao.DB.Scopes(FileInfoDB(ctx))
	if qp.Pid != 0 {
		db = db.Where("pid=?", qp.Pid)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	// db = db.Scopes(DataScope("sys_file_info", jwtauth.FromUserId(ctx)))
	// if err := db.Error; err != nil {
	// 	return nil, paginator.Info{}, err
	// }

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// Get 获取
func (cFileInfo) Get(ctx context.Context, id int) (item FileInfo, err error) {
	err = dao.DB.Scopes(FileInfoDB(ctx)).
		First(&item, "id=?", id).Error
	return
}

// Create 创建
func (cFileInfo) Create(ctx context.Context, item FileInfo) (FileInfo, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(FileInfoDB(ctx)).Create(&item).Error
	return item, err
}

// Update 更新
func (sf cFileInfo) Update(ctx context.Context, id int, up FileInfo) error {
	item, err := sf.Get(ctx, id)
	if err != nil {
		return err
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	return dao.DB.Scopes(FileInfoDB(ctx)).
		Model(&item).Updates(&up).Error
}

// Delete 删除
func (cFileInfo) Delete(ctx context.Context, id int) error {
	return dao.DB.Scopes(FileInfoDB(ctx)).
		Delete(&FileInfo{}, "id=?", id).Error
}

// BatchDelete 批量删除
func (sf cFileInfo) BatchDelete(ctx context.Context, ids []int) error {
	switch len(ids) {
	case 0:
		return nil
	case 1:
		return sf.Delete(ctx, ids[0])
	default:
		return dao.DB.Scopes(FileInfoDB(ctx)).
			Delete(&FileInfo{}, "id in (?)", ids).Error
	}
}
