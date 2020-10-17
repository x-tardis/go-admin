package models

import (
	"errors"
	_ "time"

	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/pkg/deployed"
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

func DeptTreeList(items []SysDept) []SysDept {
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

func DeptLabelTreeList(items []SysDept) []DeptLabel {
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

func (e *SysDept) Create() (SysDept, error) {
	var doc SysDept
	result := deployed.DB.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	deptPath := "/" + cast.ToString(e.DeptId)
	if int(e.ParentId) != 0 {
		var deptP SysDept
		deployed.DB.Table(e.TableName()).Where("dept_id = ?", e.ParentId).First(&deptP)
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0" + deptPath
	}
	var mp = map[string]string{}
	mp["deptPath"] = deptPath
	if err := deployed.DB.Table(e.TableName()).Where("dept_id = ?", e.DeptId).Updates(mp).Error; err != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	doc.DeptPath = deptPath
	return doc, nil
}

func (e *SysDept) Get() (SysDept, error) {
	var doc SysDept

	table := deployed.DB.Table(e.TableName())
	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}
	if e.DeptName != "" {
		table = table.Where("dept_name = ?", e.DeptName)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *SysDept) GetList() ([]SysDept, error) {
	var doc []SysDept

	table := deployed.DB.Table(e.TableName())
	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}
	if e.DeptName != "" {
		table = table.Where("dept_name = ?", e.DeptName)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.Order("sort").Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *SysDept) GetPage(bl bool) ([]SysDept, error) {
	var doc []SysDept

	table := deployed.DB.Table(e.TableName())
	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}
	if e.DeptName != "" {
		table = table.Where("dept_name = ?", e.DeptName)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}
	if e.DeptPath != "" {
		table = table.Where("deptPath like %?%", e.DeptPath)
	}
	if bl {
		// 数据权限控制
		dataPermission := new(DataPermission)
		dataPermission.UserId = cast.ToInt(e.DataScope)
		tableper, err := dataPermission.GetDataScope("sys_dept", table)
		if err != nil {
			return nil, err
		}
		table = tableper
	}

	if err := table.Order("sort").Find(&doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

func (e *SysDept) SetDept(bl bool) ([]SysDept, error) {
	list, err := e.GetPage(bl)
	if err != nil {
		return nil, err
	}
	return DeptTreeList(list), nil
}

func (e *SysDept) Update(id int) (update SysDept, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("dept_id = ?", id).First(&update).Error; err != nil {
		return
	}

	deptPath := "/" + cast.ToString(e.DeptId)
	if int(e.ParentId) != 0 {
		var deptP SysDept
		deployed.DB.Table(e.TableName()).Where("dept_id = ?", e.ParentId).First(&deptP)
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0" + deptPath
	}
	e.DeptPath = deptPath

	if e.DeptPath != "" && e.DeptPath != update.DeptPath {
		return update, errors.New("上级部门不允许修改！")
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据

	if err = deployed.DB.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}

	return
}

func (e *SysDept) Delete(id int) (success bool, err error) {
	user := SysUser{}
	user.DeptId = id
	userlist, err := user.GetList()
	if err != nil {
		return false, err
	}
	if len(userlist) <= 0 {
		return false, errors.New("当前部门存在用户，不能删除！")
	}

	tx := deployed.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		success = false
		return
	}

	if err = tx.Table(e.TableName()).Where("dept_id = ?", id).Delete(&SysDept{}).Error; err != nil {
		success = false
		tx.Rollback()
		return
	}
	if err = tx.Commit().Error; err != nil {
		success = false
		return
	}
	success = true

	return
}

func (dept *SysDept) SetDeptLabel() (m []DeptLabel, err error) {
	deptList, err := dept.GetList()
	if err != nil {
		return nil, err
	}
	return DeptLabelTreeList(deptList), nil
}
