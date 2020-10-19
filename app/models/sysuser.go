package models

import (
	"errors"
	"strings"

	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

// User
type User struct {
	// key
	IdentityKey string
	// 用户名
	UserName  string
	FirstName string
	LastName  string
	// 角色
	Role string
}

type UserName struct {
	Username string `gorm:"size:64" json:"username"`
}

type PassWord struct {
	Password string `gorm:"size:128" json:"password"`
}

type LoginM struct {
	UserName
	PassWord
}

type SysUserId struct {
	UserId int `gorm:"primary_key;AUTO_INCREMENT"  json:"userId"` // 编码
}

type SysUserB struct {
	NickName string `gorm:"size:128" json:"nickName"` // 昵称
	Phone    string `gorm:"size:11" json:"phone"`     // 手机号
	RoleId   int    `gorm:"" json:"roleId"`           // 角色编码
	Salt     string `gorm:"size:255" json:"salt"`     //盐
	Avatar   string `gorm:"size:255" json:"avatar"`   //头像
	Sex      string `gorm:"size:255" json:"sex"`      //性别
	Email    string `gorm:"size:128" json:"email"`    //邮箱
	DeptId   int    `gorm:"" json:"deptId"`           //部门编码
	PostId   int    `gorm:"" json:"postId"`           //职位编码
	CreateBy string `gorm:"size:128" json:"createBy"` //
	UpdateBy string `gorm:"size:128" json:"updateBy"` //
	Remark   string `gorm:"size:255" json:"remark"`   //备注
	Status   string `gorm:"size:4;" json:"status"`
	Model

	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
}

type SysUser struct {
	SysUserId
	LoginM
	SysUserB
}

func (SysUser) TableName() string {
	return "sys_user"
}

func UserDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(SysUser{})
	}
}

type SysUserPwd struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type SysUserPage struct {
	SysUserId
	SysUserB
	LoginM
	DeptName string `gorm:"-" json:"deptName"`
}

type SysUserView struct {
	SysUserId
	SysUserB
	LoginM
	RoleName string `gorm:"column:role_name"  json:"role_name"`
}

// 获取用户数据
func (e *SysUser) Get() (SysUserView SysUserView, err error) {

	table := deployed.DB.Table(e.TableName()).Select([]string{"sys_user.*", "sys_role.role_name"})
	table = table.Joins("left join sys_role on sys_user.role_id=sys_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.First(&SysUserView).Error; err != nil {
		return
	}

	SysUserView.Password = ""
	return
}

func (e *SysUser) GetUserInfo() (SysUserView SysUserView, err error) {
	table := deployed.DB.Table(e.TableName()).Select([]string{"sys_user.*", "sys_role.role_name"})
	table = table.Joins("left join sys_role on sys_user.role_id=sys_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.First(&SysUserView).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) GetList() (SysUserView []SysUserView, err error) {
	table := deployed.DB.Table(e.TableName()).Select([]string{"sys_user.*", "sys_role.role_name"})
	table = table.Joins("left join sys_role on sys_user.role_id=sys_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.Find(&SysUserView).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) GetPage(param paginator.Param) ([]SysUserPage, paginator.Info, error) {
	var doc []SysUserPage
	table := deployed.DB.Select("sys_user.*,sys_dept.dept_name").Table(e.TableName())
	table = table.Joins("left join sys_dept on sys_dept.dept_id = sys_user.dept_id")

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}
	if e.Status != "" {
		table = table.Where("sys_user.status = ?", e.Status)
	}

	if e.Phone != "" {
		table = table.Where("sys_user.phone = ?", e.Phone)
	}

	if e.DeptId != 0 {
		table = table.Where("sys_user.dept_id in (select dept_id from sys_dept where dept_path like ? )", "%"+cast.ToString(e.DeptId)+"%")
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(DataPermission)
	dataPermission.UserId = cast.ToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, paginator.Info{}, err
	}
	ifc, err := iorm.QueryPages(table, param, &doc)
	if err != nil {
		return nil, ifc, err
	}
	return doc, ifc, nil
}

//加密
func (e *SysUser) Encrypt() error {
	if e.Password == "" {
		return nil
	}

	hash, err := deployed.Verify.Hash(e.Password, "")
	if err != nil {
		return err
	}
	e.Password = hash
	return nil
}

//添加
func (e SysUser) Insert() (id int, err error) {
	if err = e.Encrypt(); err != nil {
		return
	}

	// check 用户名
	var count int64
	deployed.DB.Table(e.TableName()).Where("username = ?", e.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}

	//添加数据
	if err = deployed.DB.Table(e.TableName()).Create(&e).Error; err != nil {
		return
	}
	id = e.UserId
	return
}

//修改
func (e *SysUser) Update(id int) (update SysUser, err error) {
	if e.Password != "" {
		if err = e.Encrypt(); err != nil {
			return
		}
	}
	if err = deployed.DB.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}
	if e.RoleId == 0 {
		e.RoleId = update.RoleId
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = deployed.DB.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) BatchDelete(id []int) (Result bool, err error) {
	if err = deployed.DB.Table(e.TableName()).Where("user_id in (?)", id).Delete(&SysUser{}).Error; err != nil {
		return
	}
	Result = true
	return
}

func (e *SysUser) SetPwd(pwd SysUserPwd) (bool, error) {
	user, err := e.GetUserInfo()
	if err != nil {
		return false, errors.New("获取用户数据失败(代码202)")
	}
	err = deployed.Verify.Compare(pwd.OldPassword, "", user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "hashedPassword is not the hash of the given password") {
			return false, errors.New("密码错误(代码202)")
		}
		return false, err
	}
	e.Password = pwd.NewPassword
	_, err = e.Update(e.UserId)
	if err != nil {
		return false, errors.New("更新密码失败(代码202)")
	}
	return true, nil
}

func (e *SysUser) GetByUserId(tx *gorm.DB, id interface{}) error {
	return tx.First(e, id).Error
}
