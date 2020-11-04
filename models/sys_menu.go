package models

import (
	"context"
	"errors"

	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// 目录类型
const (
	MenuTypeToc  = "toc" // 菜单      M
	MenuTypeMenu = "mem" // 菜单项    C
	// MenuTypeLay  = "lay" // 布局 lay
	MenuTypeBtn = "btn" // 按钮       F
	MenuTypeIfc = "ifc" // 接口       A
)

// Menu菜单
// MenuType = toc 菜单权限, 菜单功能,里面仅含有菜单项
// MenuType = men 菜单项权限, 菜单项功能,为实际页面
//
// MenuType == btn 按钮权限, 主要用于按钮权限控制,由permission设置
//      title, paths, permission, visible
//      标题, 路径树, 权限标识, 排序, 显示
// MenuType == ifc接口权限,主要用于角色权限路由method,path的设置.
//      title, path, paths, action, sort, visible
//      标题, 路径, 路径树, 方法, 排序, 显示
type Menu struct {
	MenuId     int    `json:"menuId" gorm:"primary_key;AUTO_INCREMENT"` // 主键
	MenuName   string `json:"menuName" gorm:"size:128;"`                // 名称
	Title      string `json:"title" gorm:"size:128;"`                   // 标题
	Icon       string `json:"icon" gorm:"size:128;"`                    // 图标
	Path       string `json:"path" gorm:"size:128;"`                    // 路径
	Paths      string `json:"paths" gorm:"size:128;"`                   // 路径树
	MenuType   string `json:"menuType" gorm:"size:3;"`                  // 类型
	Action     string `json:"action" gorm:"size:16;"`                   // 请求方式,仅<接口>使用
	Permission string `json:"permission" gorm:"size:255;"`              // 权限标识,仅在<菜单项,按钮>使用
	ParentId   int    `json:"parentId" gorm:"size:11;"`                 // 父级主键
	NoCache    bool   `json:"noCache" gorm:"size:8;"`                   // 不缓存
	Breadcrumb string `json:"breadcrumb" gorm:"size:255;"`              // 面包屑
	Component  string `json:"component" gorm:"size:255;"`               // 组件路径,仅在<菜单,菜单项>使用
	Sort       int    `json:"sort"`                                     // 排序
	Visible    string `json:"visible" gorm:"size:1;"`                   // 显示/隐藏
	IsFrame    string `json:"isFrame" gorm:"size:1;DEFAULT:0;"`         // 是否外链
	Creator    string `json:"creator" gorm:"size:128;"`                 // 创建者
	Updator    string `json:"updator" gorm:"size:128;"`                 // 更新者
	Model

	Children []Menu `json:"children" gorm:"-"`

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
}

// TableName implement schema.Tabler interface
func (Menu) TableName() string {
	return "sys_menu"
}

func (sf *Menu) updatePaths(ctx context.Context) (err error) {
	curPath := "/" + cast.ToString(sf.MenuId)
	if sf.ParentId == 0 {
		sf.Paths = "/0" + curPath
	} else {
		parentMenu, err := CMenu.Get(ctx, sf.ParentId)
		if err != nil || parentMenu.Paths == "" {
			return errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
		}
		sf.Paths = parentMenu.Paths + curPath
	}
	return dao.DB.Scopes(MenuDB(ctx)).
		Where("menu_id=?", sf.MenuId).Update("paths", sf.Paths).Error
}

// MenuDB menu db scopes
func MenuDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Menu{})
	}
}

// MenuQueryParam 查询参数
type MenuQueryParam struct {
	Title    string `form:"title"`
	MenuName string `form:"menuName"`
	Path     string `form:"path"`
	Action   string `form:"action"`
	MenuType string `form:"menuType"`
	Visible  string `form:"visible"`
}

// MenuTitleLabel title树
type MenuTitleLabel struct {
	Id       int              `json:"id"`
	Label    string           `json:"label"`
	Children []MenuTitleLabel `json:"children"`
}

type cMenu struct{}

// CMenu 实例
var CMenu = cMenu{}

// toMenuTree 目录树
func toMenuTree(items []Menu) []Menu {
	tree := make([]Menu, 0)
	for _, itm := range items {
		if itm.ParentId == 0 {
			tree = append(tree, deepChildrenMenu(items, itm))
		}
	}
	return tree
}

// 获得递归子目录
func deepChildrenMenu(items []Menu, item Menu) Menu {
	item.Children = make([]Menu, 0)
	for _, itm := range items {
		if item.MenuId == itm.ParentId {
			if itm.MenuType == MenuTypeToc || itm.MenuType == MenuTypeMenu {
				itm = deepChildrenMenu(items, itm)
			}
			item.Children = append(item.Children, itm)
		}
	}
	return item
}

// 转换目录title Label树
func toMenuTitleLabelTree(items []Menu) []MenuTitleLabel {
	tree := make([]MenuTitleLabel, 0)
	for _, itm := range items {
		if itm.ParentId == 0 {
			tree = append(tree, deepChildrenMenuTitleLabel(items, MenuTitleLabel{Id: itm.MenuId, Label: itm.Title}))
		}
	}
	return tree
}

// 递归子目录title label
func deepChildrenMenuTitleLabel(items []Menu, item MenuTitleLabel) MenuTitleLabel {
	item.Children = make([]MenuTitleLabel, 0)
	for _, itm := range items {
		if item.Id == itm.ParentId {
			label := MenuTitleLabel{Id: itm.MenuId, Label: itm.Title}
			if itm.MenuType == MenuTypeToc || itm.MenuType == MenuTypeMenu {
				label = deepChildrenMenuTitleLabel(items, label)
			}
			item.Children = append(item.Children, label)
		}
	}
	return item
}

// QueryTree 获取目录树
func (sf cMenu) QueryTree(ctx context.Context, qp MenuQueryParam) ([]Menu, error) {
	items, err := sf.QueryPage(ctx, qp)
	if err != nil {
		return nil, err
	}
	return toMenuTree(items), nil
}

// QueryTreeWithRoleName 获取角色的目录树
func (sf cMenu) QueryTreeWithRoleName(ctx context.Context) ([]Menu, error) {
	roleName := jwtauth.FromRoleKey(ctx)
	items, err := sf.QueryWithRoleName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	return toMenuTree(items), nil
}

// QueryWithRoleName 通过role name 查询
func (cMenu) QueryWithRoleName(ctx context.Context, roleName string) (items []Menu, err error) {
	err = dao.DB.Scopes(MenuDB(ctx)).
		Select("sys_menu.*").
		Joins("left join sys_role_menu on sys_role_menu.menu_id=sys_menu.menu_id").
		Where("sys_role_menu.role_name=? and menu_type in (? , ?)", roleName, MenuTypeToc, MenuTypeMenu).
		Order("sort").Find(&items).Error
	return
}

// QueryTitleLabelTree 获取title 树
func (sf cMenu) QueryTitleLabelTree(ctx context.Context) ([]MenuTitleLabel, error) {
	items, err := sf.Query(ctx, MenuQueryParam{})
	if err != nil {
		return nil, err
	}
	return toMenuTitleLabelTree(items), nil
}

// Query 查询
func (cMenu) Query(ctx context.Context, qp MenuQueryParam) (items []Menu, err error) {
	db := dao.DB.Scopes(MenuDB(ctx))
	if qp.MenuName != "" {
		db = db.Where("menu_name=?", qp.MenuName)
	}
	if qp.Path != "" {
		db = db.Where("path=?", qp.Path)
	}
	if qp.Action != "" {
		db = db.Where("action=?", qp.Action)
	}
	if qp.MenuType != "" {
		db = db.Where("menu_type=?", qp.MenuType)
	}
	err = db.Order("sort").Find(&items).Error
	return
}

// QueryPage 查询,分查
func (cMenu) QueryPage(ctx context.Context, qp MenuQueryParam) (items []Menu, err error) {
	db := dao.DB.Scopes(MenuDB(ctx))
	if qp.MenuName != "" {
		db = db.Where("menu_name=?", qp.MenuName)
	}
	if qp.Title != "" {
		db = db.Where("title=?", qp.Title)
	}
	if qp.Visible != "" {
		db = db.Where("visible=?", qp.Visible)
	}
	if qp.MenuType != "" {
		db = db.Where("menu_type=?", qp.MenuType)
	}

	// 数据权限控制
	db = db.Scopes(DataScope(Menu{}, jwtauth.FromUserId(ctx)))
	if err := db.Error; err != nil {
		return nil, err
	}

	err = db.Order("sort").Find(&items).Error
	return
}

// Get 获取
func (cMenu) Get(ctx context.Context, id int) (item Menu, err error) {
	err = dao.DB.Scopes(MenuDB(ctx)).
		First(&item, "menu_id=?", id).Error
	return
}

// Create 创建目录
func (cMenu) Create(ctx context.Context, item Menu) (Menu, error) {
	err := trans.Exec(ctx, dao.DB, func(ctx context.Context) error {
		item.Creator = jwtauth.FromUserIdStr(ctx)
		err := dao.DB.Scopes(MenuDB(ctx)).Create(&item).Error
		if err != nil {
			return err
		}
		return item.updatePaths(ctx)
	})
	return item, err
}

// Update 更新
func (sf cMenu) Update(ctx context.Context, id int, up Menu) (Menu, error) {
	err := trans.Exec(ctx, dao.DB, func(ctx context.Context) error {
		oldItem, err := sf.Get(ctx, id)
		if err != nil {
			return err
		}

		up.Updator = jwtauth.FromUserIdStr(ctx)
		if err := dao.DB.Scopes(MenuDB(ctx)).Model(&oldItem).Updates(&up).Error; err != nil {
			return err
		}
		return up.updatePaths(ctx)
	})
	return up, err
}

// Delete 删除
func (sf cMenu) Delete(ctx context.Context, id int) error {
	return trans.Exec(ctx, dao.DB, func(ctx context.Context) error {
		item, err := sf.Get(ctx, id)
		if err != nil {
			return err
		}

		// 查看是否有子项
		if item.MenuType == MenuTypeToc || item.MenuType == MenuTypeMenu {
			var count int64
			dao.DB.Scopes(MenuDB(ctx)).
				Where("parent_id=?", item.MenuId).Count(&count)
			if count > 0 {
				return errors.New("有子项存在,不可删除")
			}
		}
		// 删除 sys_role_menu
		if err = CRoleMenu.DeleteWithMenu(ctx, id); err != nil {
			return err
		}
		return dao.DB.Scopes(MenuDB(ctx)).
			Delete(&Menu{}, "menu_id=?", id).Error
	})
}
