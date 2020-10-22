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

type FileInfo struct {
	Id      int    `json:"id"`                                // id
	Type    string `json:"type" gorm:"type:varchar(255);"`    // 文件类型
	Name    string `json:"name" gorm:"type:varchar(255);"`    // 文件名称
	Size    string `json:"size" gorm:"type:int(11);"`         // 文件大小
	PId     int    `json:"pId" gorm:"type:int(11);"`          // 目录id
	Source  string `json:"source" gorm:"type:varchar(255);"`  // 文件源
	Url     string `json:"url" gorm:"type:varchar(255);"`     // 文件路径
	FullUrl string `json:"fullUrl" gorm:"type:varchar(255);"` // 文件全路径
	Creator string `json:"creator" gorm:"type:varchar(128);"` // 创建人
	Updator string `json:"updator" gorm:"type:varchar(128);"` // 编辑人
	Model

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (FileInfo) TableName() string {
	return "sys_file_info"
}

func FileInfoDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(FileInfo{})
	}
}

type FileInfoQueryParam struct {
	PId int `form:"pId"`
	paginator.Param
}

type cFileInfo struct{}

var CFileInfo = new(cFileInfo)

// 获取SysFileInfo带分页
func (cFileInfo) QueryPage(ctx context.Context, qp FileInfoQueryParam) ([]FileInfo, paginator.Info, error) {
	var err error
	var items []FileInfo

	db := dao.DB.Scopes(FileInfoDB(ctx))

	if qp.PId != 0 {
		db = db.Where("p_id=?", qp.PId)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	// dataPermission := new(DataPermission)
	// dataPermission.UserId = jwtauth.FromUserId(ctx)
	// db, err = dataPermission.GetDataScope(e.TableName(), db)
	// if err != nil {
	// 	return nil, paginator.Info{}, err
	// }

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// 获取SysFileInfo
func (cFileInfo) Get(ctx context.Context, id int) (item FileInfo, err error) {
	err = dao.DB.Scopes(FileInfoDB(ctx)).
		Where("id=?", id).First(&item).Error
	return
}

// 创建SysFileInfo
func (cFileInfo) Create(ctx context.Context, item FileInfo) (FileInfo, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(FileInfoDB(ctx)).Create(&item).Error
	return item, err
}

// 更新SysFileInfo
func (cFileInfo) Update(ctx context.Context, id int, up FileInfo) (item FileInfo, err error) {
	up.Updator = jwtauth.FromUserIdStr(ctx)
	if err = dao.DB.Scopes(FileInfoDB(ctx)).
		Where("id=?", id).First(&item).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	err = dao.DB.Scopes(FileInfoDB(ctx)).
		Model(&item).Updates(&up).Error
	return
}

// 删除SysFileInfo
func (cFileInfo) Delete(ctx context.Context, id int) error {
	return dao.DB.Scopes(FileInfoDB(ctx)).
		Where("id=?", id).Delete(&FileInfo{}).Error
}

//批量删除
func (cFileInfo) BatchDelete(ctx context.Context, ids []int) error {
	return dao.DB.Scopes(FileInfoDB(ctx)).
		Where("id in (?)", ids).Delete(&FileInfo{}).Error
}
