package models

import (
	"context"

	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// RoleMenu role menu 关系表
type RoleMenu struct {
	RoleId   int    `gorm:""`
	MenuId   int    `gorm:""`
	RoleName string `gorm:"size:128"`
}

// TableName implement schema.Tabler interface
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

// RoleMenuDB role mene db scope
func RoleMenuDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(RoleMenu{})
	}
}

type cRoleMenu struct{}

// CRoleMenu 实例
var CRoleMenu = cRoleMenu{}

func (cRoleMenu) Get(ctx context.Context, roleId int) (items []RoleMenu, err error) {
	err = dao.DB.Scopes(RoleMenuDB(ctx)).
		Where("role_id=?", roleId).Find(&items).Error
	return
}

func (cRoleMenu) GetPermissionWithRoleId(ctx context.Context) ([]string, error) {
	var items []Menu

	roleId := jwtauth.FromRoleId(ctx)
	err := dao.DB.Select("sys_menu.permission").
		Table("sys_menu").
		Joins("left join sys_role_menu on sys_menu.menu_id = sys_role_menu.menu_id").
		Where("role_id = ?", roleId).
		Where("sys_menu.menu_type in(? , ?)", MenuTypeMenu, MenuTypeBtn).
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

func (cRoleMenu) GetIdsWithRoleName(ctx context.Context) (items []MenuPath, err error) {
	roleName := jwtauth.FromRoleName(ctx)
	err = dao.DB.Select("sys_menu.path").Table("sys_role_menu").
		Joins("left join sys_role on sys_role.role_id=sys_role_menu.role_id").
		Joins("left join sys_menu on sys_menu.id=sys_role_menu.menu_id").
		Where("sys_role.role_name=? and sys_menu.type=1", roleName).
		Find(&items).Error
	return
}

// BatchCreate 批量创建
func (sf cRoleMenu) BatchCreate(ctx context.Context, roleId int, menuId []int) error {
	return trans.Exec(ctx, dao.DB, func(ctx context.Context) error {
		// 获取角色的标识
		role, err := CRole.Get(ctx, roleId)
		if err != nil {
			return err
		}

		// 获取目录列表
		menus, err := CMenu.BatchGet(ctx, menuId)
		if err != nil {
			return err
		}

		roleMenus := make([]RoleMenu, 0, len(menus))
		casbinRules := make([]CasbinRule, 0, len(menus))
		for _, menu := range menus {
			roleMenus = append(roleMenus, RoleMenu{role.RoleId, menu.MenuId, role.RoleKey})
			if menu.MenuType == MenuTypeIfc {
				casbinRules = append(casbinRules, CasbinRule{
					PType: "p",
					V0:    role.RoleKey,
					V1:    menu.Path,
					V2:    menu.Method,
				})
			}
		}

		// 执行批量插入 sys_role_menu
		_, err = sf.batchCreate(ctx, roleMenus)
		if err != nil {
			return err
		}
		// 执行批量插入 sys_casbin_rule
		if len(casbinRules) > 0 {
			_, err := CCasbinRule.BatchCreate(ctx, casbinRules)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteWithRole 通过role id删除
func (cRoleMenu) DeleteWithRole(ctx context.Context, roleId int) error {
	return dao.DB.Scopes(RoleMenuDB(ctx)).
		Delete(&RoleMenu{}, "role_id=?", roleId).Error
}

// DeleteWithMenu 通过menu id删除
func (cRoleMenu) DeleteWithMenu(ctx context.Context, menuId int) error {
	return dao.DB.Scopes(RoleMenuDB(ctx)).
		Delete(&RoleMenu{}, "menu_id=?", menuId).Error
}

func (cRoleMenu) batchCreate(ctx context.Context, items []RoleMenu) ([]RoleMenu, error) {
	err := dao.DB.Scopes(RoleMenuDB(ctx)).Create(&items).Error
	return items, err
}

// GetMenuTreeOption 获取目录树和角色已选的目录id列表
func (cRoleMenu) GetMenuTreeOption(ctx context.Context, roleId int) ([]MenuTitleLabel, []int, error) {
	tree, err := CMenu.QueryTitleLabelTree(ctx)
	if err != nil {
		return nil, nil, err
	}
	menuIds := make([]int, 0)
	if roleId != 0 {
		menuIds, err = CRole.GetMenuIds(ctx, roleId)
	}
	return tree, menuIds, err
}
