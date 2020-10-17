package models

import (
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

type SysRole struct {
	RoleId    int    `json:"roleId" gorm:"primary_key;AUTO_INCREMENT"` // 角色编码
	RoleName  string `json:"roleName" gorm:"size:128;"`                // 角色名称
	Status    string `json:"status" gorm:"size:4;"`                    //
	RoleKey   string `json:"roleKey" gorm:"size:128;"`                 //角色代码
	RoleSort  int    `json:"roleSort" gorm:""`                         //角色排序
	Flag      string `json:"flag" gorm:"size:128;"`                    //
	CreateBy  string `json:"createBy" gorm:"size:128;"`                //
	UpdateBy  string `json:"updateBy" gorm:"size:128;"`                //
	Remark    string `json:"remark" gorm:"size:255;"`                  //备注
	Admin     bool   `json:"admin" gorm:"size:4;"`
	DataScope string `json:"dataScope" gorm:"size:128;"`
	Model

	Params  string `json:"params" gorm:"-"`
	MenuIds []int  `json:"menuIds" gorm:"-"`
	DeptIds []int  `json:"deptIds" gorm:"-"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

type MenuIdList struct {
	MenuId int `json:"menuId"`
}

func (role *SysRole) GetPage(param paginator.Param) ([]SysRole, paginator.Info, error) {
	var doc []SysRole

	table := deployed.DB.Table("sys_role")
	if role.RoleId != 0 {
		table = table.Where("role_id = ?", role.RoleId)
	}
	if role.RoleName != "" {
		table = table.Where("role_name = ?", role.RoleName)
	}
	if role.Status != "" {
		table = table.Where("status = ?", role.Status)
	}
	if role.RoleKey != "" {
		table = table.Where("role_key = ?", role.RoleKey)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId = cast.ToInt(role.DataScope)
	table, err := dataPermission.GetDataScope("sys_role", table)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	ifc, err := iorm.QueryPages(table.Order("role_sort"), param, &doc)
	return doc, ifc, err
}

func (role *SysRole) Get() (SysRole SysRole, err error) {
	table := deployed.DB.Table("sys_role")
	if role.RoleId != 0 {
		table = table.Where("role_id = ?", role.RoleId)
	}
	if role.RoleName != "" {
		table = table.Where("role_name = ?", role.RoleName)
	}
	if err = table.First(&SysRole).Error; err != nil {
		return
	}
	return
}

func (role *SysRole) GetList() (SysRole []SysRole, err error) {
	table := deployed.DB.Table("sys_role")
	if role.RoleId != 0 {
		table = table.Where("role_id = ?", role.RoleId)
	}
	if role.RoleName != "" {
		table = table.Where("role_name = ?", role.RoleName)
	}
	if err = table.Order("role_sort").Find(&SysRole).Error; err != nil {
		return
	}

	return
}

// 获取角色对应的菜单ids
func (*SysRole) GetRoleMenuId(roleId int) ([]int, error) {
	menuIds := make([]int, 0)
	menuList := make([]MenuIdList, 0)
	if err := deployed.DB.Table("sys_role_menu").Select("sys_role_menu.menu_id").Where("role_id = ? ", roleId).
		Where(" sys_role_menu.menu_id not in(select sys_menu.parent_id from sys_role_menu LEFT JOIN sys_menu on sys_menu.menu_id=sys_role_menu.menu_id where role_id =?  and parent_id is not null)", roleId).
		Find(&menuList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(menuList); i++ {
		menuIds = append(menuIds, menuList[i].MenuId)
	}
	return menuIds, nil
}

func (role *SysRole) Insert() (id int, err error) {
	var i int64
	deployed.DB.Table(role.TableName()).Where("role_name=? or role_key = ?", role.RoleName, role.RoleKey).Count(&i)
	if i > 0 {
		return 0, errors.New("角色名称或者角色标识已经存在！")
	}
	role.UpdateBy = ""
	result := deployed.DB.Table(role.TableName()).Create(&role)
	if result.Error != nil {
		err = result.Error
		return
	}
	id = role.RoleId
	return
}

type DeptIdList struct {
	DeptId int `json:"DeptId"`
}

func (role *SysRole) GetRoleDeptId() ([]int, error) {
	deptIds := make([]int, 0)
	deptList := make([]DeptIdList, 0)
	if err := deployed.DB.Table("sys_role_dept").Select("sys_role_dept.dept_id").
		Joins("LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id").Where("role_id = ? ", role.RoleId).
		Where(" sys_role_dept.dept_id not in(select sys_dept.parent_id from sys_role_dept LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id where role_id =? )", role.RoleId).
		Find(&deptList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(deptList); i++ {
		deptIds = append(deptIds, deptList[i].DeptId)
	}

	return deptIds, nil
}

//修改
func (role *SysRole) Update(id int) (update SysRole, err error) {
	if err = deployed.DB.Table(role.TableName()).First(&update, id).Error; err != nil {
		return
	}

	if role.RoleName != "" && role.RoleName != update.RoleName {
		return update, errors.New("角色名称不允许修改！")
	}

	if role.RoleKey != "" && role.RoleKey != update.RoleKey {
		return update, errors.New("角色标识不允许修改！")
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = deployed.DB.Table(role.TableName()).Model(&update).Updates(&role).Error; err != nil {
		return
	}
	return
}

//批量删除
func (role *SysRole) BatchDelete(id []int) (Result bool, err error) {
	tx := deployed.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}
	// 查询角色
	var roles []SysRole
	if err := tx.Table("sys_role").Where("role_id in (?)", id).Find(&roles).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	var count int64
	if err := tx.Table("sys_user").Where("role_id in (?)", id).Count(&count).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	if count > 0 {
		tx.Rollback()
		return false, errors.New("存在绑定用户，请解绑后重试")
	}

	// 删除角色
	if err = tx.Table(role.TableName()).Where("role_id in (?)", id).Unscoped().Delete(&SysRole{}).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 删除角色菜单
	if err := tx.Table("sys_role_menu").Where("role_id in (?)", id).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 删除casbin配置
	for i := 0; i < len(roles); i++ {
		if err := tx.Table("sys_casbin_rule").Where("v0 in (?)", roles[0].RoleKey).Delete(&CasbinRule{}).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	return true, nil
}
