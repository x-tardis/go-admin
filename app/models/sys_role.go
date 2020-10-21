package models

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type Role struct {
	RoleId    int    `json:"roleId" gorm:"primary_key;AUTO_INCREMENT"` // 角色编码
	RoleName  string `json:"roleName" gorm:"size:128;"`                // 角色名称
	RoleKey   string `json:"roleKey" gorm:"size:128;"`                 // 角色代码
	Status    string `json:"status" gorm:"size:4;"`                    // 状态
	Sort      int    `json:"sort" gorm:""`                             // 角色排序
	Flag      string `json:"flag" gorm:"size:128;"`                    //
	Remark    string `json:"remark" gorm:"size:255;"`                  // 备注
	Admin     bool   `json:"admin" gorm:"size:4;"`
	DataScope string `json:"dataScope" gorm:"size:128;"`
	Creator   string `json:"creator" gorm:"size:128;"` // 创建者
	Updator   string `json:"updator" gorm:"size:128;"` // 更新者
	Model

	MenuIds []int `json:"menuIds" gorm:"-"`
	DeptIds []int `json:"deptIds" gorm:"-"`

	Params string `json:"params" gorm:"-"`
}

func (Role) TableName() string {
	return "sys_role"
}

func RoleDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(Role{})
	}
}

type RoleQueryParam struct {
	RoleName string `form:"roleName"`
	RoleKey  string `form:"roleKey"`
	Status   string `form:"status"`
	paginator.Param
}

type cRole struct{}

var CRole = new(cRole)

func (cRole) Query(_ context.Context) (item []Role, err error) {
	err = deployed.DB.Scopes(RoleDB()).
		Order("sort").
		Find(&item).Error
	return
}

func (cRole) QueryPage(ctx context.Context, qp RoleQueryParam) ([]Role, paginator.Info, error) {
	var err error
	var items []Role

	db := deployed.DB.Scopes(RoleDB())
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
	db, err = GetDataScope("sys_role", jwtauth.FromUserId(ctx), db)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

func (cRole) GetWithName(_ context.Context, name string) (item Role, err error) {
	err = deployed.DB.Scopes(RoleDB()).
		Where("role_name=?", name).
		First(&item).Error
	return
}

func (cRole) Get(_ context.Context, id int) (item Role, err error) {
	err = deployed.DB.Scopes(RoleDB()).
		Where("role_id=?", id).
		First(&item).Error
	return
}

func (cRole) Create(ctx context.Context, item Role) (Role, error) {
	var i int64

	deployed.DB.Table(item.TableName()).
		Where("role_name=? or role_key = ?", item.RoleName, item.RoleKey).Count(&i)
	if i > 0 {
		return item, errors.New("角色名称或者角色标识已经存在！")
	}

	item.Creator = jwtauth.FromUserIdStr(ctx)
	item.Updator = ""
	err := deployed.DB.Scopes(RoleDB()).Create(&item).Error
	return item, err
}

// 修改
func (cRole) Update(ctx context.Context, id int, up Role) (item Role, err error) {
	up.Updator = jwtauth.FromUserIdStr(ctx)
	if err = deployed.DB.Scopes(RoleDB()).First(&item, id).Error; err != nil {
		return
	}

	if up.RoleName != "" && up.RoleName != item.RoleName {
		return item, errors.New("角色名称不允许修改！")
	}

	if up.RoleKey != "" && up.RoleKey != item.RoleKey {
		return item, errors.New("角色标识不允许修改！")
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = deployed.DB.Scopes(RoleDB()).Model(&item).Updates(&up).Error
	return
}

// 批量删除
func (cRole) BatchDelete(_ context.Context, id []int) error {
	tx := deployed.DB.Begin()

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
	if err := tx.Scopes(RoleDB()).Where("role_id in (?)", id).Find(&roles).Error; err != nil {
		tx.Rollback()
		return err
	}

	var count int64
	if err := tx.Scopes(UserDB()).Where("role_id in (?)", id).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count > 0 {
		tx.Rollback()
		return errors.New("存在绑定用户，请解绑后重试")
	}

	// 删除角色
	if err := tx.Scopes(RoleDB()).Where("role_id in (?)", id).Unscoped().Delete(&Role{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除角色菜单
	if err := tx.Scopes(RoleMenuDB()).Where("role_id in (?)", id).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除casbin配置
	for i := 0; i < len(roles); i++ {
		if err := tx.Scopes(CasbinRuleDB()).Where("v0 in (?)", roles[0].RoleKey).Delete(&CasbinRule{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

type MenuIdList struct {
	MenuId int `json:"menuId"`
}

// 获取角色对应的菜单ids
func (cRole) GetMenuIds(roleId int) ([]int, error) {
	var menuList []MenuIdList

	if err := deployed.DB.Scopes(RoleMenuDB()).
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

type DeptIdList struct {
	DeptId int `json:"DeptId"`
}

func (cRole) GetDeptIds(_ context.Context, roleId int) ([]int, error) {
	var deptList []DeptIdList

	if err := deployed.DB.Scopes(RoleDeptDB()).
		Select("sys_role_dept.dept_id").
		Joins("LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id").
		Where("role_id = ? ", roleId).
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
