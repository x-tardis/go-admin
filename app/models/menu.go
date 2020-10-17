package models

import (
	"errors"

	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/pkg/deployed"
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
	CreateBy   string `json:"createBy" gorm:"size:128;"`
	UpdateBy   string `json:"updateBy" gorm:"size:128;"`
	IsFrame    string `json:"isFrame" gorm:"size:1;DEFAULT:0;"`
	DataScope  string `json:"dataScope" gorm:"-"`
	Params     string `json:"params" gorm:"-"`
	RoleId     int    `gorm:"-"`
	Children   []Menu `json:"children" gorm:"-"`
	IsSelect   bool   `json:"is_select" gorm:"-"`
	Model
}

func (Menu) TableName() string {
	return "sys_menu"
}

// MenuTreeList 目录树
func MenuTreeList(items []Menu) []Menu {
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

type MenuLabel struct {
	Id       int         `json:"id" gorm:"-"`
	Label    string      `json:"label" gorm:"-"`
	Children []MenuLabel `json:"children" gorm:"-"`
}

// MenuLabelTreeList 目录Lable树
func MenuLabelTreeList(items []Menu) []MenuLabel {
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

type Menus struct {
	MenuId     int    `json:"menuId" gorm:"column:menu_id;primary_key;"`
	MenuName   string `json:"menuName" gorm:"column:menu_name"`
	Title      string `json:"title" gorm:"column:title"`
	Icon       string `json:"icon" gorm:"column:icon"`
	Path       string `json:"path" gorm:"column:path"`
	MenuType   string `json:"menuType" gorm:"column:menu_type"`
	Action     string `json:"action" gorm:"column:action"`
	Permission string `json:"permission" gorm:"column:permission"`
	ParentId   int    `json:"parentId" gorm:"column:parent_id"`
	NoCache    bool   `json:"noCache" gorm:"column:no_cache"`
	Breadcrumb string `json:"breadcrumb" gorm:"column:breadcrumb"`
	Component  string `json:"component" gorm:"column:component"`
	Sort       int    `json:"sort" gorm:"column:sort"`

	Visible  string `json:"visible" gorm:"column:visible"`
	Children []Menu `json:"children" gorm:"-"`

	CreateBy  string `json:"createBy" gorm:"column:create_by"`
	UpdateBy  string `json:"updateBy" gorm:"column:update_by"`
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
	Model
}

func (Menus) TableName() string {
	return "sys_menu"
}

type MenuRole struct {
	Menus
	IsSelect bool `json:"is_select" gorm:"-"`
}

func (e *Menu) GetByMenuId(id int) (Menu Menu, err error) {
	table := deployed.DB.Table(e.TableName())
	err = table.Where("menu_id = ?", id).Find(&Menu).Error
	return
}

func (e *Menu) SetMenu() (m []Menu, err error) {
	menuList, err := e.GetPage()
	if err != nil {
		return nil, err
	}

	return MenuTreeList(menuList), nil
}

func (e *Menu) SetMenuLabel() (m []MenuLabel, err error) {
	menuList, err := e.Get()
	if err != nil {
		return nil, err
	}
	return MenuLabelTreeList(menuList), nil
}

func (e *Menu) SetMenuRole(rolename string) ([]Menu, error) {
	menulist, err := e.GetByRoleName(rolename)
	if err != nil {
		return nil, err
	}
	return MenuTreeList(menulist), nil
}

func (e *MenuRole) Get() (Menus []MenuRole, err error) {
	table := deployed.DB.Table(e.TableName())
	if e.MenuName != "" {
		table = table.Where("menu_name = ?", e.MenuName)
	}
	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	return
}

func (e *Menu) GetByRoleName(rolename string) (Menus []Menu, err error) {
	table := deployed.DB.Table(e.TableName()).Select("sys_menu.*").Joins("left join sys_role_menu on sys_role_menu.menu_id=sys_menu.menu_id")
	table = table.Where("sys_role_menu.role_name=? and menu_type in ('M','C')", rolename)
	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	return
}

func (e *Menu) Get() (Menus []Menu, err error) {
	table := deployed.DB.Table(e.TableName())
	if e.MenuName != "" {
		table = table.Where("menu_name = ?", e.MenuName)
	}
	if e.Path != "" {
		table = table.Where("path = ?", e.Path)
	}
	if e.Action != "" {
		table = table.Where("action = ?", e.Action)
	}
	if e.MenuType != "" {
		table = table.Where("menu_type = ?", e.MenuType)
	}

	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	return
}

func (e *Menu) GetPage() (Menus []Menu, err error) {
	table := deployed.DB.Table(e.TableName())
	if e.MenuName != "" {
		table = table.Where("menu_name = ?", e.MenuName)
	}
	if e.Title != "" {
		table = table.Where("title = ?", e.Title)
	}
	if e.Visible != "" {
		table = table.Where("visible = ?", e.Visible)
	}
	if e.MenuType != "" {
		table = table.Where("menu_type = ?", e.MenuType)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId = cast.ToInt(e.DataScope)
	table, err = dataPermission.GetDataScope("sys_menu", table)
	if err != nil {
		return nil, err
	}
	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	return
}

func (e *Menu) Create() (id int, err error) {
	result := deployed.DB.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err = result.Error
		return
	}
	err = InitPaths(e)
	if err != nil {
		return
	}
	id = e.MenuId
	return
}

func InitPaths(menu *Menu) (err error) {
	parentMenu := new(Menu)
	if menu.ParentId != 0 {
		deployed.DB.Table("sys_menu").Where("menu_id = ?", menu.ParentId).First(parentMenu)
		if parentMenu.Paths == "" {
			return errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
		}
		menu.Paths = parentMenu.Paths + "/" + cast.ToString(menu.MenuId)
	} else {
		menu.Paths = "/0/" + cast.ToString(menu.MenuId)
	}
	return deployed.DB.Table("sys_menu").Where("menu_id = ?", menu.MenuId).Update("paths", menu.Paths).Error
}

func (e *Menu) Update(id int) (update Menu, err error) {
	if err = deployed.DB.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = deployed.DB.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	err = InitPaths(e)
	return
}

func (e *Menu) Delete(id int) error {
	return deployed.DB.Table(e.TableName()).Where("menu_id = ?", id).Delete(&Menu{}).Error
}
