package models

import (
	"context"
	"errors"
	_ "time"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type Dept struct {
	DeptId   int    `json:"deptId" gorm:"primary_key;auto_increment;"` // 部门编码
	ParentId int    `json:"parentId" gorm:""`                          // 上级部门
	DeptPath string `json:"deptPath" gorm:"size:255;"`                 //
	DeptName string `json:"deptName"  gorm:"size:128;"`                // 部门名称
	Sort     int    `json:"sort" gorm:""`                              // 排序
	Leader   string `json:"leader" gorm:"size:128;"`                   // 负责人
	Phone    string `json:"phone" gorm:"size:11;"`                     // 手机
	Email    string `json:"email" gorm:"size:64;"`                     // 邮箱
	Status   string `json:"status" gorm:"size:4;"`                     // 状态
	Creator  string `json:"creator" gorm:"size:64;"`
	Updator  string `json:"updator" gorm:"size:64;"`
	Model

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
	Children  []Dept `json:"children" gorm:"-"`
}

func (Dept) TableName() string {
	return "sys_dept"
}

func DeptDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(Dept{})
	}
}

type DeptLabel struct {
	Id       int         `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLabel `gorm:"-" json:"children"`
}

type DeptQueryParam struct {
	DeptId   int    `form:"deptId"`
	DeptName string `form:"deptName"`
	DeptPath string `form:"deptPath"`
	Status   string `form:"status"`
}

type cDept struct{}

var CDept = new(cDept)

func toDeptTree(items []Dept) []Dept {
	tree := make([]Dept, 0)
	for _, itm := range items {
		if itm.ParentId == 0 {
			tree = append(tree, deepChildrenDept(items, itm))
		}
	}
	return tree
}

func deepChildrenDept(items []Dept, item Dept) Dept {
	item.Children = make([]Dept, 0)
	for _, itm := range items {
		if item.DeptId == itm.ParentId {
			item.Children = append(item.Children, deepChildrenDept(items, itm))
		}
	}
	return item
}

func toDeptLabelTree(items []Dept) []DeptLabel {
	tree := make([]DeptLabel, 0)
	for _, itm := range items {
		if itm.ParentId == 0 {
			tree = append(tree, deepChildrenDeptLabel(items, DeptLabel{
				itm.DeptId,
				itm.DeptName,
				make([]DeptLabel, 0),
			}))
		}
	}
	return tree
}

func deepChildrenDeptLabel(items []Dept, dept DeptLabel) DeptLabel {
	dept.Children = make([]DeptLabel, 0)
	for _, itm := range items {
		if dept.Id == itm.ParentId {
			dept.Children = append(dept.Children, deepChildrenDeptLabel(items, DeptLabel{
				itm.DeptId,
				itm.DeptName,
				make([]DeptLabel, 0),
			}))
		}
	}
	return dept
}

func (sf cDept) QueryLabelTree(ctx context.Context) ([]DeptLabel, error) {
	items, err := sf.Query(ctx)
	if err != nil {
		return nil, err
	}
	return toDeptLabelTree(items), nil
}

func (sf cDept) QueryTree(ctx context.Context, qp DeptQueryParam, bl bool) ([]Dept, error) {
	list, err := sf.QueryPage(ctx, qp, bl)
	if err != nil {
		return nil, err
	}
	return toDeptTree(list), nil
}

func (cDept) Query(_ context.Context) (items []Dept, err error) {
	err = deployed.DB.Scopes(DeptDB()).
		Order("sort").Find(&items).Error
	return
}

func (cDept) QueryPage(ctx context.Context, qp DeptQueryParam, bl bool) (items []Dept, err error) {
	db := deployed.DB.Scopes(DeptDB())
	if qp.DeptId != 0 {
		db = db.Where("dept_id=?", qp.DeptId)
	}
	if qp.DeptName != "" {
		db = db.Where("dept_name=?", qp.DeptName)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}
	if qp.DeptPath != "" {
		db = db.Where("deptPath like %?%", qp.DeptPath)
	}

	if bl {
		// 数据权限控制
		dataPermission := new(DataPermission)
		dataPermission.UserId = jwtauth.FromUserId(ctx)
		db, err = dataPermission.GetDataScope("sys_dept", db)
		if err != nil {
			return nil, err
		}
	}

	err = db.Order("sort").Find(&items).Error
	return items, err
}

func (cDept) Get(_ context.Context, id int) (item Dept, err error) {
	err = deployed.DB.Scopes(DeptDB()).
		Where("dept_id=?", id).First(&item).Error
	return
}

func (cDept) Create(ctx context.Context, item Dept) (Dept, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := deployed.DB.Scopes(DeptDB()).Create(&item).Error
	if err != nil {
		return item, err
	}

	deptPath := "/" + cast.ToString(item.DeptId)
	if item.ParentId == 0 {
		deptPath = "/0" + deptPath
	} else {
		var parentDept Dept
		deployed.DB.Scopes(DeptDB()).
			Where("dept_id=?", item.ParentId).First(&parentDept)
		deptPath = parentDept.DeptPath + deptPath
	}

	item.DeptPath = deptPath
	err = deployed.DB.Scopes(DeptDB()).
		Where("dept_id=?", item.DeptId).
		Updates(map[string]interface{}{"dept_path": deptPath}).Error
	return item, err
}

func (cDept) Update(ctx context.Context, id int, up Dept) (item Dept, err error) {
	up.Updator = jwtauth.FromUserIdStr(ctx)
	if err = deployed.DB.Scopes(DeptDB()).
		Where("dept_id=?", id).First(&item).Error; err != nil {
		return
	}

	deptPath := "/" + cast.ToString(id)
	if up.ParentId == 0 {
		deptPath = "/0" + deptPath

	} else {
		var parentDept Dept
		deployed.DB.Scopes(DeptDB()).
			Where("dept_id=?", up.ParentId).First(&parentDept)
		deptPath = parentDept.DeptPath + deptPath
	}
	up.DeptPath = deptPath

	if up.DeptPath != "" && up.DeptPath != item.DeptPath {
		return item, errors.New("上级部门不允许修改！")
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = deployed.DB.Scopes(DeptDB()).
		Model(&item).Updates(&up).Error
	return
}

func (cDept) Delete(_ context.Context, id int) error {
	userList, err := new(cUser).GetWithDeptId(id)
	if err != nil {
		return err
	}
	if len(userList) > 0 {
		return errors.New("当前部门存在用户，不能删除！")
	}

	tx := deployed.DB.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Scopes(DeptDB()).Where("dept_id=?", id).Delete(&Dept{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
