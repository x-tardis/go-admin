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

// Dept 部门
type Dept struct {
	DeptId   int    `json:"deptId" gorm:"primary_key;auto_increment;"` // 主键
	ParentId int    `json:"parentId" gorm:""`                          // 上级主键
	DeptPath string `json:"deptPath" gorm:"size:255;"`                 // 路径树
	DeptName string `json:"deptName"  gorm:"size:128;"`                // 名称
	Sort     int    `json:"sort" gorm:""`                              // 排序
	Leader   string `json:"leader" gorm:"size:128;"`                   // 负责人
	Phone    string `json:"phone" gorm:"size:11;"`                     // 负责人联系手机
	Email    string `json:"email" gorm:"size:64;"`                     // 负责人联系邮箱
	Status   string `json:"status" gorm:"size:4;"`                     // 状态
	Creator  string `json:"creator" gorm:"size:64;"`                   // 创建者
	Updator  string `json:"updator" gorm:"size:64;"`                   // 更新者
	Model

	Children []Dept `json:"children" gorm:"-"` // 子列表

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
}

// TableName 实现gorm.Tabler接口
func (Dept) TableName() string {
	return "sys_dept"
}

// DeptDB scope dept model
func DeptDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(Dept{})
	}
}

// DeptLabel dept label
type DeptLabel struct {
	Id       int         `json:"id"`
	Label    string      `json:"label"`
	Children []DeptLabel `json:"children"`
}

// DeptQueryParam 查询参数
type DeptQueryParam struct {
	DeptId   int    `form:"deptId"`
	DeptName string `form:"deptName"`
	DeptPath string `form:"deptPath"`
	Status   string `form:"status"`
}

type cDept struct{}

// CDept 实例
var CDept = cDept{}

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
				nil,
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

// QueryLabelTree query label tree
func (sf cDept) QueryLabelTree(ctx context.Context) ([]DeptLabel, error) {
	items, err := sf.Query(ctx)
	if err != nil {
		return nil, err
	}
	return toDeptLabelTree(items), nil
}

// QueryTree query tree
func (sf cDept) QueryTree(ctx context.Context, qp DeptQueryParam, bl bool) ([]Dept, error) {
	list, err := sf.QueryPage(ctx, qp, bl)
	if err != nil {
		return nil, err
	}
	return toDeptTree(list), nil
}

// Query 查询部门列表
func (cDept) Query(_ context.Context) (items []Dept, err error) {
	err = deployed.DB.Scopes(DeptDB()).
		Order("sort").Find(&items).Error
	return
}

// QueryPage 查询部门列表,分页
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

	// 数据权限控制
	if bl {
		db, err = GetDataScope("sys_dept", jwtauth.FromUserId(ctx), db)
		if err != nil {
			return nil, err
		}
	}

	err = db.Order("sort").Find(&items).Error
	return items, err
}

// Get 通过Id获取部门
func (cDept) Get(_ context.Context, id int) (item Dept, err error) {
	err = deployed.DB.Scopes(DeptDB()).
		Where("dept_id=?", id).First(&item).Error
	return
}

// Create 创建部门
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

// Update 更新部门信息
func (cDept) Update(ctx context.Context, id int, up Dept) (item Dept, err error) {
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

	up.Updator = jwtauth.FromUserIdStr(ctx)
	err = deployed.DB.Scopes(DeptDB()).
		Model(&item).Updates(&up).Error
	return
}

// Delete 删除部门
func (cDept) Delete(_ context.Context, id int) error {
	userList, err := CUser.GetWithDeptId(id)
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
