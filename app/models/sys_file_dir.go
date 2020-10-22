package models

import (
	"context"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/app/models/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type FileDir struct {
	Id      int    `json:"id"`
	Label   string `json:"label" gorm:"type:varchar(255);"`   // 名称
	PId     int    `json:"pId" gorm:"type:int(11);"`          // 父id
	Sort    int    `json:"sort" gorm:""`                      // 排序
	Path    string `json:"path" gorm:"size:255;"`             //
	Creator string `json:"creator" gorm:"type:varchar(128);"` // 创建人
	Updator string `json:"updator" gorm:"type:varchar(128);"` // 编辑人
	Model

	Children  []FileDir `json:"children" gorm:"-"`
	DataScope string    `json:"dataScope" gorm:"-"`
	Params    string    `json:"params"  gorm:"-"`
}

func (FileDir) TableName() string {
	return "sys_file_dir"
}

func FileDirDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(FileDir{})
	}
}

type FileDirQueryParam struct {
	Id    int    `form:"id"`
	Label string `form:"label"`
	PId   int    `form:"pId"`
}

type cFileDir struct{}

var CFileDir = new(cFileDir)

func toFileDirTree(items []FileDir) []FileDir {
	tree := make([]FileDir, 0)
	for _, itm := range items {
		if itm.PId == 0 {
			tree = append(tree, deepChildrenFileDir(items, itm))
		}
	}
	return tree
}

func deepChildrenFileDir(items []FileDir, item FileDir) FileDir {
	item.Children = make([]FileDir, 0)
	for _, itm := range items {
		if item.Id == itm.PId {
			item.Children = append(item.Children, deepChildrenFileDir(items, itm))
		}
	}
	return item
}

func (sf cFileDir) QueryTree(ctx context.Context, qp FileDirQueryParam) ([]FileDir, error) {
	item, err := sf.Query(ctx, qp)
	if err != nil {
		return nil, err
	}
	return toFileDirTree(item), nil
}

// 获取SysFileDir带分页
func (cFileDir) Query(ctx context.Context, qp FileDirQueryParam) ([]FileDir, error) {
	var err error
	var items []FileDir

	db := dao.DB.Scopes(FileDirDB())
	if qp.Id != 0 {
		db = db.Where("id=?", qp.Id)
	}
	if qp.Label != "" {
		db = db.Where("label=?", qp.Label)
	}
	if qp.PId != 0 {
		db = db.Where("p_id=?", qp.PId)
	}

	// // 数据权限控制(如果不需要数据权限请将此处去掉)
	// dataPermission := new(DataPermission)
	// dataPermission.UserId = jwtauth.FromUserId(ctx)
	// db, err = dataPermission.GetDataScope(e.TableName(), db)
	// if err != nil {
	// 	return nil, 0, err
	// }

	err = db.Find(&items).Error
	return items, err
}

// 获取SysFileDir
func (cFileDir) Get(_ context.Context, id int) (item FileDir, err error) {
	err = dao.DB.Scopes(FileDirDB()).
		Where("id=?", id).First(&item).Error
	return
}

// 创建SysFileDir
func (cFileDir) Create(ctx context.Context, item FileDir) (FileDir, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(FileDirDB()).Create(&item).Error
	if err != nil {
		return item, err
	}

	path := "/" + cast.ToString(item.Id)
	if item.PId == 0 {
		path = "/0" + path
	} else {
		var deptP FileDir
		dao.DB.Scopes(FileDirDB()).Where("id = ?", item.PId).First(&deptP)
		path = deptP.Path + path
	}

	err = dao.DB.Scopes(FileDirDB()).
		Where("id = ?", item.Id).
		Updates(map[string]interface{}{"path": path}).Error
	item.Path = path
	return item, err
}

// 更新SysFileDir
func (cFileDir) Update(ctx context.Context, id int, up FileDir) (item FileDir, err error) {
	up.Updator = jwtauth.FromUserIdStr(ctx)
	if err = dao.DB.Scopes(FileDirDB()).
		Where("id=?", id).First(&item).Error; err != nil {
		return
	}

	path := "/" + cast.ToString(up.Id)
	if up.Id == 0 {
		path = "/0" + path
	} else {
		var deptP FileDir
		dao.DB.Scopes(FileDirDB()).Where("id = ?", up.Id).First(&deptP)
		path = deptP.Path + path
	}
	up.Path = path

	// if up.Path != "" && up.Path != item.Path {
	//	return item, errors.New("上级不允许修改！")
	// }

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = dao.DB.Scopes(FileDirDB()).
		Model(&item).Updates(&up).Error
	return
}

// 删除SysFileDir
func (cFileDir) Delete(_ context.Context, id int) error {
	return dao.DB.Scopes(FileDirDB()).
		Where("id=?", id).Delete(&FileDir{}).Error
}

// 批量删除
func (cFileDir) BatchDelete(_ context.Context, ids []int) error {
	return dao.DB.Scopes(FileDirDB()).
		Where("id in (?)", ids).Delete(&FileDir{}).Error
}
