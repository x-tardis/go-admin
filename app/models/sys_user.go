package models

import (
	"context"
	"errors"

	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
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
	Salt     string `gorm:"size:255" json:"salt"`     // 盐
	Avatar   string `gorm:"size:255" json:"avatar"`   // 头像
	Sex      string `gorm:"size:255" json:"sex"`      // 性别
	Email    string `gorm:"size:128" json:"email"`    // 邮箱
	DeptId   int    `gorm:"" json:"deptId"`           // 部门编码
	PostId   int    `gorm:"" json:"postId"`           // 职位编码
	Remark   string `gorm:"size:255" json:"remark"`   // 备注
	Status   string `gorm:"size:4;" json:"status"`
	Creator  string `gorm:"size:128" json:"creator"` //
	Updator  string `gorm:"size:128" json:"updator"` //
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

type UpdateUserPwd struct {
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

type UserQueryParam struct {
	Username string `form:"username"`
	Status   string `form:"status"`
	Phone    string `form:"phone"`
	PostId   int    `form:"postId"`
	DeptId   int    `form:"deptId"`
	paginator.Param
}

type CallUser struct{}

// 获取用户数据
func (sf CallUser) Get(ctx context.Context, id int) (SysUserView, error) {
	item, err := sf.getUserInfo(ctx, id)
	if err != nil {
		return SysUserView{}, err
	}
	item.Password = ""
	return item, err
}

// 获取用户数据
func (sf CallUser) GetUserInfo(ctx context.Context) (SysUserView, error) {
	return sf.Get(ctx, jwtauth.FromUserId(ctx))
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

func (CallUser) QueryPage(ctx context.Context, qp UserQueryParam) ([]SysUserPage, paginator.Info, error) {
	var items []SysUserPage

	db := deployed.DB.Scopes(UserDB()).
		Select("sys_user.*,sys_dept.dept_name").
		Joins("left join sys_dept on sys_dept.dept_id = sys_user.dept_id")

	if qp.Username != "" {
		db = db.Where("username=?", qp.Username)
	}
	if qp.Status != "" {
		db = db.Where("sys_user.status=?", qp.Status)
	}
	if qp.Phone != "" {
		db = db.Where("sys_user.phone=?", qp.Phone)
	}
	if qp.DeptId != 0 {
		db = db.Where("sys_user.dept_id in (select dept_id from sys_dept where dept_path like ? )", "%"+cast.ToString(qp.DeptId)+"%")
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(DataPermission)
	dataPermission.UserId = jwtauth.FromUserId(ctx)
	db, err := dataPermission.GetDataScope("sys_user", db)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// 加密
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

// 添加
func (CallUser) Create(ctx context.Context, item SysUser) (SysUser, error) {
	var count int64
	var err error

	// check 用户名
	deployed.DB.Scopes(UserDB()).Where("username=?", item.Username).Count(&count)
	if count > 0 {
		return item, errors.New("账户已存在！")
	}

	item.Password, err = deployed.Verify.Hash(item.Password, "")
	if err != nil {
		return item, err
	}
	item.Creator = jwtauth.FromUserIdStr(ctx)
	// 添加数据
	err = deployed.DB.Table(item.TableName()).Create(&item).Error
	return item, err
}

// 修改
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

	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	if err = deployed.DB.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (CallUser) BatchDelete(id []int) error {
	return deployed.DB.Scopes(UserDB()).
		Where("user_id in (?)", id).Delete(&SysUser{}).Error
}

func (sf CallUser) UpdatePassword(ctx context.Context, pwd UpdateUserPwd) error {
	item, err := sf.getUserInfo(ctx, jwtauth.FromUserId(ctx))
	if err != nil {
		return errors.New("获取用户数据失败(代码202)")
	}

	// 校验旧密码 和 新密加签
	err = deployed.Verify.Compare(pwd.OldPassword, "", item.Password)
	if err != nil {
		return err
	}
	pass, err := deployed.Verify.Hash(pwd.NewPassword, "")
	if err != nil {
		return err
	}

	err = deployed.DB.Scopes(UserDB()).
		Where("user_id=?", item.UserId).Update("password", pass).Error
	if err != nil {
		return errors.New("更新密码失败(代码202)")
	}
	return nil
}

func (CallUser) getUserInfo(_ context.Context, id int) (item SysUserView, err error) {
	err = deployed.DB.Scopes(UserDB()).
		Select([]string{"sys_user.*", "sys_role.role_name"}).
		Joins("left join sys_role on sys_user.role_id=sys_role.role_id").
		Where("user_id=?", id).First(&item).Error
	return
}
