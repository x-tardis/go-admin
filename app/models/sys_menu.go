package models

import (
	"context"
	"errors"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type Menu struct {
	MenuId     int    `json:"menuId" gorm:"primary_key;AUTO_INCREMENT"`
	MenuName   string `json:"menuName" gorm:"size:128;"`
	Title      string `json:"title" gorm:"size:128;"`
	Icon       string `json:"icon" gorm:"size:128;"`
	Path       string `json:"path" gorm:"size:128;"`
	Paths      string `json:"paths" gorm:"size:128;"`
	MenuType   string `json:"menuType" gorm:"size:1;"`
	Action     string `json:"action" gorm:"size:16;"`
	Permission string `json:"permission" gorm:"size:255;"`
	ParentId   int    `json:"parentId" gorm:"size:11;"`
	NoCache    bool   `json:"noCache" gorm:"size:8;"`
	Breadcrumb string `json:"breadcrumb" gorm:"size:255;"`
	Component  string `json:"component" gorm:"size:255;"`
	Sort       int    `json:"sort" gorm:"size:4;"`
	Visible    string `json:"visible" gorm:"size:1;"`
	IsFrame    string `json:"isFrame" gorm:"size:1;DEFAULT:0;"`
	RoleId     int    `gorm:"-"`
	Children   []Menu `json:"children" gorm:"-"`
	IsSelect   bool   `json:"is_select" gorm:"-"`
	Creator    string `json:"creator" gorm:"size:128;"`
	Updator    string `json:"updator" gorm:"size:128;"`
	Model

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
}

func (Menu) TableName() string {
	return "sys_menu"
}

func MenuDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(Menu{})
	}
}

type MenuQueryParam struct {
	MenuName string `form:"menuName"`
	Path     string `form:"path"`
	Action   string `form:"action"`
	MenuType string `form:"menuType"`
	Visible  string `form:"visible"`
	Title    string `form:"title"`
}
type MenuLabel struct {
	Id       int         `json:"id" gorm:"-"`
	Label    string      `json:"label" gorm:"-"`
	Children []MenuLabel `json:"children" gorm:"-"`
}

type cMenu struct{}

var CMenu = new(cMenu)

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

// deepChildrenMenu 获得递归子目录
func deepChildrenMenu(items []Menu, item Menu) Menu {
	item.Children = make([]Menu, 0)
	for _, itm := range items {
		if item.MenuId == itm.ParentId {
			if itm.MenuType != "F" {
				itm = deepChildrenMenu(items, itm)
			}
			item.Children = append(item.Children, itm)
		}
	}
	return item
}

// toMenuLabelTree 目录Label树
func toMenuLabelTree(items []Menu) []MenuLabel {
	tree := make([]MenuLabel, 0)
	for _, itm := range items {
		if itm.ParentId == 0 {
			lab := MenuLabel{
				itm.MenuId,
				itm.Title,
				make([]MenuLabel, 0),
			}
			tree = append(tree, deepChildrenMenuLabel(items, lab))
		}
	}
	return tree
}

// deepChildrenMenuLabel 获得递归子目录Lable
func deepChildrenMenuLabel(items []Menu, item MenuLabel) MenuLabel {
	for _, itm := range items {
		if item.Id == itm.ParentId {
			mi := MenuLabel{
				itm.MenuId,
				itm.Title,
				make([]MenuLabel, 0),
			}
			if itm.MenuType != "F" {
				mi = deepChildrenMenuLabel(items, mi)
			}
			item.Children = append(item.Children, mi)
		}
	}
	return item
}

//
// type Menus struct {
// 	MenuId     int    `json:"menuId" gorm:"column:menu_id;primary_key;"`
// 	MenuName   string `json:"menuName" gorm:"column:menu_name"`
// 	Title      string `json:"title" gorm:"column:title"`
// 	Icon       string `json:"icon" gorm:"column:icon"`
// 	Path       string `json:"path" gorm:"column:path"`
// 	MenuType   string `json:"menuType" gorm:"column:menu_type"`
// 	Action     string `json:"action" gorm:"column:action"`
// 	Permission string `json:"permission" gorm:"column:permission"`
// 	ParentId   int    `json:"parentId" gorm:"column:parent_id"`
// 	NoCache    bool   `json:"noCache" gorm:"column:no_cache"`
// 	Breadcrumb string `json:"breadcrumb" gorm:"column:breadcrumb"`
// 	Component  string `json:"component" gorm:"column:component"`
// 	Sort       int    `json:"sort" gorm:"column:sort"`
//
// 	Visible  string `json:"visible" gorm:"column:visible"`
// 	Children []Menu `json:"children" gorm:"-"`
//
// 	Creator  string `json:"creator" gorm:"column:creator"`
// 	Updator  string `json:"updator" gorm:"column:updator"`
// 	DataScope string `json:"dataScope" gorm:"-"`
// 	Params    string `json:"params" gorm:"-"`
// 	Model
// }
//
// func (Menus) TableName() string {
// 	return "sys_menu"
// }

// type MenuRole struct {
// 	Menus
// 	IsSelect bool `json:"is_select" gorm:"-"`
// }
//
// func (*MenuRole) Get(menuName string) (items []MenuRole, err error) {
// 	db := deployed.DB.Scopes(MenuDB())
// 	if menuName != "" {
// 		db = db.Where("menu_name = ?", menuName)
// 	}
// 	err = db.Order("sort").Find(&items).Error
// 	return
// }

func (sf cMenu) QueryTree(ctx context.Context, qp MenuQueryParam) (m []Menu, err error) {
	items, err := sf.QueryPage(ctx, qp)
	if err != nil {
		return nil, err
	}
	return toMenuTree(items), nil
}

func (sf cMenu) QueryTreeWithRoleName(ctx context.Context, roleName string) ([]Menu, error) {
	items, err := sf.QueryWithRoleName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	return toMenuTree(items), nil
}

func (cMenu) QueryWithRoleName(_ context.Context, roleName string) (items []Menu, err error) {
	err = deployed.DB.Scopes(MenuDB()).
		Select("sys_menu.*").
		Joins("left join sys_role_menu on sys_role_menu.menu_id=sys_menu.menu_id").
		Where("sys_role_menu.role_name=? and menu_type in ('M','C')", roleName).
		Order("sort").Find(&items).Error
	return
}

func (sf cMenu) QueryLabelTree(ctx context.Context) (m []MenuLabel, err error) {
	items, err := sf.Query(ctx, MenuQueryParam{})
	if err != nil {
		return nil, err
	}
	return toMenuLabelTree(items), nil
}

func (cMenu) Query(_ context.Context, qp MenuQueryParam) (items []Menu, err error) {
	db := deployed.DB.Scopes(MenuDB())
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

func (cMenu) QueryPage(ctx context.Context, qp MenuQueryParam) (items []Menu, err error) {
	db := deployed.DB.Scopes(MenuDB())
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
	dataPermission := new(DataPermission)
	dataPermission.UserId = jwtauth.FromUserId(ctx)
	db, err = dataPermission.GetDataScope("sys_menu", db)
	if err != nil {
		return nil, err
	}

	err = db.Order("sort").Find(&items).Error
	return
}

func (cMenu) Get(ctx context.Context, id int) (item Menu, err error) {
	err = deployed.DB.Scopes(MenuDB()).Where("menu_id=?", id).Find(&item).Error
	return
}
func (cMenu) Create(ctx context.Context, item Menu) (Menu, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := deployed.DB.Scopes(MenuDB()).Create(&item).Error
	if err != nil {
		return item, err
	}
	err = InitPaths(&item)
	return item, err
}

func (cMenu) Update(ctx context.Context, id int, up Menu) (item Menu, err error) {
	up.Updator = jwtauth.FromUserIdStr(ctx)
	if err = deployed.DB.Scopes(MenuDB()).First(&item, id).Error; err != nil {
		return
	}

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	if err = deployed.DB.Scopes(MenuDB()).Model(&item).Updates(&up).Error; err != nil {
		return
	}
	err = InitPaths(&up)
	return
}

func (cMenu) Delete(_ context.Context, id int) error {
	return deployed.DB.Scopes(MenuDB()).
		Where("menu_id=?", id).Delete(&Menu{}).Error
}

func InitPaths(menu *Menu) (err error) {
	parentMenu := new(Menu)
	if menu.ParentId != 0 {
		deployed.DB.Scopes(MenuDB()).
			Where("menu_id = ?", menu.ParentId).First(parentMenu)
		if parentMenu.Paths == "" {
			return errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
		}
		menu.Paths = parentMenu.Paths + "/" + cast.ToString(menu.MenuId)
	} else {
		menu.Paths = "/0/" + cast.ToString(menu.MenuId)
	}
	return deployed.DB.Scopes(MenuDB()).
		Where("menu_id = ?", menu.MenuId).Update("paths", menu.Paths).Error
}
