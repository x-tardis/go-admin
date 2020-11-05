package models

import (
	"context"

	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type FileDir struct {
	Id      int    `json:"id"`                                // 主键
	Label   string `json:"label" gorm:"type:varchar(255);"`   // 名称
	Pid     int    `json:"pid" gorm:"type:int(11);"`          // 父id
	Path    string `json:"path" gorm:"size:255;"`             // 路径树
	Sort    int    `json:"sort" gorm:""`                      // 排序
	Creator string `json:"creator" gorm:"type:varchar(128);"` // 创建者
	Updator string `json:"updator" gorm:"type:varchar(128);"` // 更新者
	Model

	Children []FileDir `json:"children" gorm:"-"`

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

// TableName implement schema.Tabler interface
func (FileDir) TableName() string {
	return "sys_file_dir"
}

// FileDirDB file dir db scopes
func FileDirDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(FileDir{})
	}
}

// FileDirQueryParam 文件查询
type FileDirQueryParam struct {
	Id    int    `form:"id"`
	Label string `form:"label"`
	Pid   int    `form:"pid"`
}

type cFileDir struct{}

// CFileDir 实例
var CFileDir = cFileDir{}

func toFileDirTree(items []FileDir) []FileDir {
	tree := make([]FileDir, 0)
	for _, itm := range items {
		if itm.Pid == 0 {
			tree = append(tree, deepChildrenFileDir(items, itm))
		}
	}
	return tree
}

func deepChildrenFileDir(items []FileDir, item FileDir) FileDir {
	item.Children = make([]FileDir, 0)
	for _, itm := range items {
		if item.Id == itm.Pid {
			item.Children = append(item.Children, deepChildrenFileDir(items, itm))
		}
	}
	return item
}

// QueryTree 查询文件路径树
func (sf cFileDir) QueryTree(ctx context.Context, qp FileDirQueryParam) ([]FileDir, error) {
	item, err := sf.Query(ctx, qp)
	if err != nil {
		return nil, err
	}
	return toFileDirTree(item), nil
}

// Query 查询FileDir带分页
func (cFileDir) Query(ctx context.Context, qp FileDirQueryParam) ([]FileDir, error) {
	var items []FileDir

	db := dao.DB.Scopes(FileDirDB(ctx))
	if qp.Id != 0 {
		db = db.Where("id=?", qp.Id)
	}
	if qp.Label != "" {
		db = db.Where("label=?", qp.Label)
	}
	if qp.Pid != 0 {
		db = db.Where("pid=?", qp.Pid)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	db = db.Scopes(DataScope(FileDir{}, jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, err
	}

	err := db.Find(&items).Error
	return items, err
}

// Get 获取SysFileDir
func (cFileDir) Get(ctx context.Context, id int) (item FileDir, err error) {
	err = dao.DB.Scopes(FileDirDB(ctx)).
		First(&item, "id=?", id).Error
	return
}

// Create 创建
func (sf cFileDir) Create(ctx context.Context, item FileDir) (FileDir, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(FileDirDB(ctx)).Create(&item).Error
	if err != nil {
		return item, err
	}

	path := "/" + cast.ToString(item.Id)
	if item.Pid == 0 {
		item.Path = "/0" + path
	} else {
		deptP, err := sf.Get(ctx, item.Pid)
		if err != nil {
			return item, err
		}
		item.Path = deptP.Path + path
	}

	err = dao.DB.Scopes(FileDirDB(ctx)).
		Where("id=?", item.Id).
		Update("path", item.Path).Error
	return item, err
}

// Update 更新
func (sf cFileDir) Update(ctx context.Context, id int, up FileDir) error {
	item, err := sf.Get(ctx, id)
	if err != nil {
		return err
	}

	path := "/" + cast.ToString(up.Id)
	if up.Pid == 0 {
		path = "/0" + path
	} else {
		deptP, err := sf.Get(ctx, up.Pid)
		if err != nil {
			return err
		}
		path = deptP.Path + path
	}
	up.Path = path

	// if up.Path != "" && up.Path != item.Path {
	//	return item, errors.New("上级不允许修改！")
	// }
	up.Updator = jwtauth.FromUserIdStr(ctx)
	return dao.DB.Scopes(FileDirDB(ctx)).
		Model(&item).Updates(&up).Error
}

// Delete 删除
func (cFileDir) Delete(ctx context.Context, id int) error {
	return dao.DB.Scopes(FileDirDB(ctx)).
		Delete(&FileDir{}, "id=?", id).Error
}

// BatchDelete 批量删除
func (sf cFileDir) BatchDelete(ctx context.Context, ids []int) error {
	switch len(ids) {
	case 0:
		return nil
	case 1:
		return sf.Delete(ctx, ids[0])
	default:
		return dao.DB.Scopes(FileDirDB(ctx)).
			Delete(&FileDir{}, "id in (?)", ids).Error
	}
}
