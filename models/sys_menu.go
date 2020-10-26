package models

import (
	"context"
	"errors"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/trans"
)

const (
	MenuTypeToc  = "toc" // 目录 M
	MenuTypeMenu = "mem" // 菜单 C
	// MenuTypeLay  = "lay" // 布局 lay
	MenuTypeBtn = "btn" // 按钮 F
	MenuTypeIfc = "ifc" // 接口 A
)

type Menu struct {
	MenuId     int    `json:"menuId" gorm:"primary_key;AUTO_INCREMENT"` // 主键
	MenuName   string `json:"menuName" gorm:"size:128;"`                // 名称
	Title      string `json:"title" gorm:"size:128;"`                   // 标题
	Icon       string `json:"icon" gorm:"size:128;"`                    // 图标
	Path       string `json:"path" gorm:"size:128;"`                    // 目录路径
	Paths      string `json:"paths" gorm:"size:128;"`                   // 路径树
	MenuType   string `json:"menuType" gorm:"size:3;"`                  // 类型
	Action     string `json:"action" gorm:"size:16;"`                   // 请求方式,仅接口使用
	Permission string `json:"permission" gorm:"size:255;"`              // 权限标识,仅在(菜单,按钮)使用
	ParentId   int    `json:"parentId" gorm:"size:11;"`                 // 父级主键
	NoCache    bool   `json:"noCache" gorm:"size:8;"`                   // 不缓存
	Breadcrumb string `json:"breadcrumb" gorm:"size:255;"`              // 面包屑
	Component  string `json:"component" gorm:"size:255;"`               // 组件路径
	Sort       int    `json:"sort" gorm:"size:4;"`                      // 排序
	Visible    string `json:"visible" gorm:"size:1;"`                   // 显示/隐藏
	IsFrame    string `json:"isFrame" gorm:"size:1;DEFAULT:0;"`         // 是否外链
	Creator    string `json:"creator" gorm:"size:128;"`                 // 创建者
	Updator    string `json:"updator" gorm:"size:128;"`                 // 更新者
	Model

	RoleId   int    `gorm:"-"`
	Children []Menu `json:"children" gorm:"-"`
	IsSelect bool   `json:"is_select" gorm:"-"`

	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
}

func (Menu) TableName() string {
	return "sys_menu"
}

func MenuDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Menu{})
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

// deepChildrenMenu 获得递归子目录
func deepChildrenMenu(items []Menu, item Menu) Menu {
	item.Children = make([]Menu, 0)
	for _, itm := range items {
		if item.MenuId == itm.ParentId {
			if itm.MenuType != MenuTypeBtn {
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
			if itm.MenuType != MenuTypeBtn {
				mi = deepChildrenMenuLabel(items, mi)
			}
			item.Children = append(item.Children, mi)
		}
	}
	return item
}

func (sf cMenu) QueryTree(ctx context.Context, qp MenuQueryParam) (m []Menu, err error) {
	items, err := sf.QueryPage(ctx, qp)
	if err != nil {
		return nil, err
	}
	return toMenuTree(items), nil
}

func (sf cMenu) QueryTreeWithRoleName(ctx context.Context) ([]Menu, error) {
	roleName := jwtauth.FromRoleKey(ctx)
	items, err := sf.QueryWithRoleName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	return toMenuTree(items), nil
}

func (cMenu) QueryWithRoleName(ctx context.Context, roleName string) (items []Menu, err error) {
	err = dao.DB.Scopes(MenuDB(ctx)).
		Select("sys_menu.*").
		Joins("left join sys_role_menu on sys_role_menu.menu_id=sys_menu.menu_id").
		Where("sys_role_menu.role_name=? and menu_type in (? , ?)", roleName, MenuTypeToc, MenuTypeMenu).
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

func (cMenu) Get(ctx context.Context, id int) (item Menu, err error) {
	err = dao.DB.Scopes(MenuDB(ctx)).Where("menu_id=?", id).Find(&item).Error
	return
}
func (cMenu) Create(ctx context.Context, item Menu) (Menu, error) {
	item.Creator = jwtauth.FromUserIdStr(ctx)
	err := dao.DB.Scopes(MenuDB(ctx)).Create(&item).Error
	if err != nil {
		return item, err
	}
	err = InitPaths(ctx, &item)
	return item, err
}

func (cMenu) Update(ctx context.Context, id int, up Menu) (item Menu, err error) {
	if err = dao.DB.Scopes(MenuDB(ctx)).First(&item, id).Error; err != nil {
		return
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	if err = dao.DB.Scopes(MenuDB(ctx)).Model(&item).Updates(&up).Error; err != nil {
		return
	}
	err = InitPaths(ctx, &up)
	return
}

func (cMenu) Delete(ctx context.Context, id int) error {
	return dao.DB.Scopes(MenuDB(ctx)).
		Where("menu_id=?", id).Delete(&Menu{}).Error
}

func InitPaths(ctx context.Context, menu *Menu) (err error) {
	parentMenu := new(Menu)
	if menu.ParentId != 0 {
		dao.DB.Scopes(MenuDB(ctx)).
			Where("menu_id = ?", menu.ParentId).First(parentMenu)
		if parentMenu.Paths == "" {
			return errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
		}
		menu.Paths = parentMenu.Paths + "/" + cast.ToString(menu.MenuId)
	} else {
		menu.Paths = "/0/" + cast.ToString(menu.MenuId)
	}
	return dao.DB.Scopes(MenuDB(ctx)).
		Where("menu_id = ?", menu.MenuId).Update("paths", menu.Paths).Error
}
