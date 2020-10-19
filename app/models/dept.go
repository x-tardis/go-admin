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

type SysDept struct {
	DeptId   int    `json:"deptId" gorm:"primary_key;auto_increment;"` // 部门编码
	ParentId int    `json:"parentId" gorm:""`                          // 上级部门
	DeptPath string `json:"deptPath" gorm:"size:255;"`                 //
	DeptName string `json:"deptName"  gorm:"size:128;"`                // 部门名称
	Sort     int    `json:"sort" gorm:""`                              // 排序
	Leader   string `json:"leader" gorm:"size:128;"`                   // 负责人
	Phone    string `json:"phone" gorm:"size:11;"`                     // 手机
	Email    string `json:"email" gorm:"size:64;"`                     // 邮箱
	Status   string `json:"status" gorm:"size:4;"`                     // 状态
	CreateBy string `json:"createBy" gorm:"size:64;"`
	UpdateBy string `json:"updateBy" gorm:"size:64;"`
	Model

	DataScope string    `json:"dataScope" gorm:"-"`
	Params    string    `json:"params" gorm:"-"`
	Children  []SysDept `json:"children" gorm:"-"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}

func DeptDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(SysDept{})
	}
}

type DeptQueryParam struct {
	DeptId   int    `form:"deptId"`
	DeptName string `form:"deptName"`
	DeptPath string `form:"deptPath"`
	Status   string `form:"status"`
}

func toDeptTree(items []SysDept) []SysDept {
	tree := make([]SysDept, 0)
	for _, itm := range items {
		if itm.ParentId == 0 {
			tree = append(tree, deepChildrenDept(items, itm))
		}
	}
	return tree
}

func deepChildrenDept(items []SysDept, item SysDept) SysDept {
	item.Children = make([]SysDept, 0)
	for _, itm := range items {
		if item.DeptId == itm.ParentId {
			item.Children = append(item.Children, deepChildrenDept(items, itm))
		}
	}
	return item
}

type DeptLabel struct {
	Id       int         `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLabel `gorm:"-" json:"children"`
}

type CallDept struct{}

func toDeptLabelTree(items []SysDept) []DeptLabel {
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

func deepChildrenDeptLabel(items []SysDept, dept DeptLabel) DeptLabel {
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

func (sf CallDept) QueryLabelTree(ctx context.Context) ([]DeptLabel, error) {
	items, err := sf.Query(ctx)
	if err != nil {
		return nil, err
	}
	return toDeptLabelTree(items), nil
}

func (sf CallDept) QueryTree(ctx context.Context, qp DeptQueryParam, bl bool) ([]SysDept, error) {
	list, err := sf.QueryPage(ctx, qp, bl)
	if err != nil {
		return nil, err
	}
	return toDeptTree(list), nil
}

func (CallDept) Query(_ context.Context) (items []SysDept, err error) {
	err = deployed.DB.Scopes(DeptDB()).
		Order("sort").Find(&items).Error
	return
}

func (CallDept) QueryPage(ctx context.Context, qp DeptQueryParam, bl bool) (items []SysDept, err error) {
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
		tableper, err := dataPermission.GetDataScope("sys_dept", db)
		if err != nil {
			return nil, err
		}
		db = tableper
	}

	err = db.Order("sort").Find(&items).Error
	return items, err
}

func (CallDept) Get(_ context.Context, id int) (item SysDept, err error) {
	err = deployed.DB.Scopes(DeptDB()).
		Where("dept_id=?", id).First(&item).Error
	return
}

func (CallDept) Create(ctx context.Context, item SysDept) (SysDept, error) {
	item.CreateBy = jwtauth.FromUserIdStr(ctx)
	err := deployed.DB.Scopes(DeptDB()).Create(&item).Error
	if err != nil {
		return item, err
	}

	deptPath := "/" + cast.ToString(item.DeptId)
	if item.ParentId == 0 {
		deptPath = "/0" + deptPath
	} else {
		var parentDept SysDept
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

func (CallDept) Update(ctx context.Context, id int, up SysDept) (item SysDept, err error) {
	up.UpdateBy = jwtauth.FromUserIdStr(ctx)
	if err = deployed.DB.Scopes(DeptDB()).
		Where("dept_id=?", id).First(&item).Error; err != nil {
		return
	}

	deptPath := "/" + cast.ToString(id)
	if up.ParentId == 0 {
		deptPath = "/0" + deptPath

	} else {
		var parentDept SysDept
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

func (CallDept) Delete(_ context.Context, id int) error {
	user := SysUser{}
	user.DeptId = id
	userlist, err := user.GetList()
	if err != nil {
		return err
	}
	if len(userlist) <= 0 {
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

	if err := tx.Scopes(DeptDB()).Where("dept_id=?", id).Delete(&SysDept{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
