package models

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/trans"
)

// SuperAdmin 超级管理员
const (
	SuperAdmin = "admin"
)

// Role 角色
type Role struct {
	RoleId    int    `json:"roleId" gorm:"primary_key;AUTO_INCREMENT"` // 主键
	RoleName  string `json:"roleName" gorm:"size:128;"`                // 名称
	RoleKey   string `json:"roleKey" gorm:"size:128;"`                 // 标识
	Status    string `json:"status" gorm:"size:4;"`                    // 状态
	Sort      int    `json:"sort" gorm:""`                             // 排序
	Flag      string `json:"flag" gorm:"size:128;"`                    // 标记(未用)
	Admin     bool   `json:"admin" gorm:"size:4;"`                     // 超级权限(未用)
	Remark    string `json:"remark" gorm:"size:255;"`                  // 备注
	DataScope string `json:"dataScope" gorm:"size:128;"`               // 数据权限
	Creator   string `json:"creator" gorm:"size:128;"`                 // 创建者
	Updator   string `json:"updator" gorm:"size:128;"`                 // 更新者
	Model

	MenuIds []int `json:"menuIds" gorm:"-"` // 角色目录ID列表
	DeptIds []int `json:"deptIds" gorm:"-"` // 角色部门ID列表

	Params string `json:"params" gorm:"-"` // (未用)
}

// TableName implement schema.Tabler interface
func (Role) TableName() string {
	return "sys_role"
}

// RoleDB role db scope
func RoleDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Role{})
	}
}

// RoleQueryParam 查询参数
type RoleQueryParam struct {
	RoleName string `form:"roleName"`
	RoleKey  string `form:"roleKey"`
	Status   string `form:"status"`
	paginator.Param
}

type cRole struct{}

// CRole role实例
var CRole = cRole{}

// Query 获取角色列表
func (cRole) Query(ctx context.Context) (items []Role, err error) {
	err = dao.DB.Scopes(RoleDB(ctx)).
		Order("sort").Find(&items).Error
	return
}

// QueryPage 分页查询角色列表
func (cRole) QueryPage(ctx context.Context, qp RoleQueryParam) ([]Role, paginator.Info, error) {
	var err error
	var items []Role

	db := dao.DB.Scopes(RoleDB(ctx))
	if qp.RoleName != "" {
		db = db.Where("role_name=?", qp.RoleName)
	}
	if qp.RoleKey != "" {
		db = db.Where("role_key=?", qp.RoleKey)
	}
	if qp.Status != "" {
		db = db.Where("status=?", qp.Status)
	}
	db = db.Order("sort")

	// 数据权限控制
	db = db.Scopes(DataScope(Role{}, jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// GetWithName 通过角色名获取角色信息
func (cRole) GetWithName(ctx context.Context, name string) (item Role, err error) {
	err = dao.DB.Scopes(RoleDB(ctx)).
		Where("role_name=?", name).First(&item).Error
	return
}

// Get 通过id获取角色信息
func (cRole) Get(ctx context.Context, id int) (item Role, err error) {
	err = dao.DB.Scopes(RoleDB(ctx)).
		Where("role_id=?", id).First(&item).Error
	return
}

// Create 创建角色
func (cRole) Create(ctx context.Context, item Role) (Role, error) {
	var i int64

	dao.DB.Table(item.TableName()).
		Where("role_name=? or role_key = ?", item.RoleName, item.RoleKey).Count(&i)
	if i > 0 {
		return item, errors.New("角色名称或者角色标识已经存在！")
	}

	item.Creator = jwtauth.FromUserIdStr(ctx)
	item.Updator = item.Creator
	err := dao.DB.Scopes(RoleDB(ctx)).Create(&item).Error
	return item, err
}

// Update 修改角色信息
func (cRole) Update(ctx context.Context, id int, up Role) (item Role, err error) {
	if err = dao.DB.Scopes(RoleDB(ctx)).First(&item, id).Error; err != nil {
		return
	}

	if item.RoleKey == SuperAdmin && up.Status == StatusDisable {
		return item, errors.New("超级用户不允许禁用")
	}
	// 角色名称与角色标识不允许修改
	if up.RoleName != "" && up.RoleName != item.RoleName {
		return item, errors.New("角色名称不允许修改！")
	}
	if up.RoleKey != "" && up.RoleKey != item.RoleKey {
		return item, errors.New("角色标识不允许修改！")
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	err = dao.DB.Scopes(RoleDB(ctx)).Model(&item).Updates(&up).Error
	return
}

// BatchDelete 批量删除
func (cRole) BatchDelete(ctx context.Context, ids []int) error {
	tx := dao.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	// 查询角色
	var roles []Role
	if err := tx.Scopes(RoleDB(ctx)).Where("role_id in (?)", ids).Find(&roles).Error; err != nil {
		tx.Rollback()
		return err
	}

	var count int64
	if err := tx.Scopes(UserDB(ctx)).Where("role_id in (?)", ids).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count > 0 {
		tx.Rollback()
		return errors.New("存在绑定用户，请解绑后重试")
	}

	// 删除角色
	if err := tx.Scopes(RoleDB(ctx)).Where("role_id in (?)", ids).Unscoped().Delete(&Role{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除角色菜单
	if err := tx.Scopes(RoleMenuDB(ctx)).Where("role_id in (?)", ids).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除casbin配置
	for i := 0; i < len(roles); i++ {
		if err := tx.Scopes(CasbinRuleDB(ctx)).Where("v0 in (?)", roles[0].RoleKey).Delete(&CasbinRule{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// MenuIdList ...
type MenuIdList struct {
	MenuId int `json:"menuId"`
}

// GetMenuIds 获取角色对应的菜单id列表
func (cRole) GetMenuIds(ctx context.Context, roleId int) ([]int, error) {
	var menuList []MenuIdList

	if err := dao.DB.Scopes(RoleMenuDB(ctx)).
		Select("sys_role_menu.menu_id").
		Where("role_id=? ", roleId).
		Where(" sys_role_menu.menu_id not in(select sys_menu.parent_id from sys_role_menu LEFT JOIN sys_menu on sys_menu.menu_id=sys_role_menu.menu_id where role_id=?  and parent_id is not null)", roleId).
		Find(&menuList).Error; err != nil {
		return nil, err
	}
	menuIds := make([]int, 0)
	for _, v := range menuList {
		menuIds = append(menuIds, v.MenuId)
	}
	return menuIds, nil
}

// DeptIdList ...
type DeptIdList struct {
	DeptId int `json:"DeptId"`
}

// GetDeptIds 获取角色对应的部门id列表
func (cRole) GetDeptIds(ctx context.Context, roleId int) ([]int, error) {
	var deptList []DeptIdList

	if err := dao.DB.Scopes(RoleDeptDB(ctx)).
		Select("sys_role_dept.dept_id").
		Joins("LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id").
		Where("role_id=? ", roleId).
		Where(" sys_role_dept.dept_id not in(select sys_dept.parent_id from sys_role_dept LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id where role_id =? )", roleId).
		Find(&deptList).Error; err != nil {
		return nil, err
	}

	deptIds := make([]int, 0)
	for _, v := range deptList {
		deptIds = append(deptIds, v.DeptId)
	}
	return deptIds, nil
}
