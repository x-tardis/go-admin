package models

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// RoleMenu role menu
type RoleMenu struct {
	RoleId   int    `gorm:""`
	MenuId   int    `gorm:""`
	RoleName string `gorm:"size:128"`
}

// TableName implement gorm.Tabler interface
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

// RoleMenuDB role mene db scope
func RoleMenuDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(RoleMenu{})
	}
}

type cRoleMenu struct{}

var CRoleMenu = new(cRoleMenu)

func (cRoleMenu) Get(_ context.Context, id int) (items []RoleMenu, err error) {
	err = dao.DB.Scopes(RoleMenuDB()).
		Where("role_id = ?", id).Find(&items).Error
	return
}

func (cRoleMenu) GetPermissionWithRoleId(ctx context.Context) ([]string, error) {
	var items []Menu

	roleId := jwtauth.FromRoleId(ctx)
	err := dao.DB.Select("sys_menu.permission").
		Table("sys_menu").
		Joins("left join sys_role_menu on sys_menu.menu_id = sys_role_menu.menu_id").
		Where("role_id = ?", roleId).
		Where("sys_menu.menu_type in('F','C')").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	list := make([]string, 0, len(items))
	for _, item := range items {
		list = append(list, item.Permission)
	}
	return list, nil
}

type MenuPath struct {
	Path string `json:"path"`
}

func (cRoleMenu) GetIDSWithRoleName(ctx context.Context) (items []MenuPath, err error) {
	roleName := jwtauth.FromRoleName(ctx)
	err = dao.DB.Select("sys_menu.path").Table("sys_role_menu").
		Joins("left join sys_role on sys_role.role_id=sys_role_menu.role_id").
		Joins("left join sys_menu on sys_menu.id=sys_role_menu.menu_id").
		Where("sys_role.role_name=? and sys_menu.type=1", roleName).
		Find(&items).Error
	return
}

func (cRoleMenu) DeleteWith(RoleId string, MenuID string) error {
	db := dao.DB.Scopes(RoleMenuDB()).Where("role_id = ?", RoleId)
	if MenuID != "" {
		db = db.Where("menu_id = ?", MenuID)
	}
	return db.Delete(&RoleMenu{}).Error
}

func (cRoleMenu) Delete(roleId int) error {
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Scopes(RoleDeptDB()).
		Where("role_id=?", roleId).Delete(&RoleDept{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Scopes(RoleMenuDB()).
		Where("role_id = ?", roleId).Delete(RoleMenuDB()).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var role Role

	err = tx.Scopes(RoleDB()).
		Where("role_id = ?", roleId).First(&role).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	sql3 := "delete from sys_casbin_rule where v0= '" + role.RoleKey + "';"
	if err := tx.Exec(sql3).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (cRoleMenu) Insert(roleId int, menuId []int) error {
	var (
		role            Role
		menu            []Menu
		casbinRuleQueue []CasbinRule // casbinRule 待插入队列
	)

	// 开始事务
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// 在事务中做一些数据库操作（从这一点使用'tx'，而不是'db'）
	err := tx.Scopes(RoleDB()).
		Where("role_id = ?", roleId).First(&role).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Scopes(MenuDB()).
		Where("menu_id in (?)", menuId).Find(&menu).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// ORM不支持批量插入所以需要拼接 sql 串
	var roleMenus []RoleMenu

	casbinRuleSql := "INSERT INTO `sys_casbin_rule`  (`p_type`,`v0`,`v1`,`v2`) VALUES "

	for _, m := range menu {
		roleMenus = append(roleMenus,
			RoleMenu{role.RoleId, m.MenuId, role.RoleKey})
		if m.MenuType == "A" {
			// 加入队列
			casbinRuleQueue = append(casbinRuleQueue,
				CasbinRule{V0: role.RoleKey, V1: m.Path, V2: m.Action})
		}
	}
	// 执行批量插入sys_role_menu
	err = tx.Scopes(RoleMenuDB()).Create(&roleMenus).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 拼装'sys_casbin_rule'批量插入SQL语句
	// TODO: casbinRuleQueue队列不为空时才会拼装，否则直接忽略不执行'for'循环
	for i, v := range casbinRuleQueue {
		casbinRuleSql += fmt.Sprintf("('p','%s','%s','%s')", v.V0, v.V1, v.V2)
		if i == len(casbinRuleQueue)-1 {
			casbinRuleSql += ";"
		} else {
			casbinRuleSql += ","
		}
	}
	// 执行批量插入sys_casbin_rule
	if len(casbinRuleQueue) > 0 {
		if err := tx.Exec(casbinRuleSql).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
