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

type User struct {
	UserId   int    `gorm:"primary_key;AUTO_INCREMENT"  json:"userId"` // 编码
	Username string `gorm:"size:64" json:"username"`                   // 用户名
	Password string `gorm:"size:128" json:"password"`                  // 密码
	NickName string `gorm:"size:128" json:"nickName"`                  // 昵称
	Phone    string `gorm:"size:11" json:"phone"`                      // 手机号
	RoleId   int    `gorm:"" json:"roleId"`                            // 角色编码
	Salt     string `gorm:"size:255" json:"salt"`                      // 盐
	Avatar   string `gorm:"size:255" json:"avatar"`                    // 头像
	Sex      string `gorm:"size:255" json:"sex"`                       // 性别
	Email    string `gorm:"size:128" json:"email"`                     // 邮箱
	DeptId   int    `gorm:"" json:"deptId"`                            // 部门编码
	PostId   int    `gorm:"" json:"postId"`                            // 职位编码
	Remark   string `gorm:"size:255" json:"remark"`                    // 备注
	Status   string `gorm:"size:4;" json:"status"`                     // 状态
	Creator  string `gorm:"size:128" json:"creator"`                   // 创建者
	Updator  string `gorm:"size:128" json:"updator"`                   // 更新者
	Model

	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
}

func (User) TableName() string {
	return "sys_user"
}

func UserDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(User{})
	}
}

type UserPage struct {
	User
	DeptName string `gorm:"-" json:"deptName"`
}

type UserView struct {
	User
	RoleName string `gorm:"column:role_name"  json:"roleName"`
}

type UserQueryParam struct {
	Username string `form:"username"`
	Status   string `form:"status"`
	Phone    string `form:"phone"`
	PostId   int    `form:"postId"`
	DeptId   int    `form:"deptId"`
	paginator.Param
}

type UpdateUserPwd struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type cUser struct{}

var CUser = new(cUser)

// 获取用户数据
func (sf cUser) Get(ctx context.Context, id int) (UserView, error) {
	item, err := sf.get(ctx, id)
	if err != nil {
		return UserView{}, err
	}
	item.Password = ""
	return item, err
}

// 获取用户数据
func (sf cUser) GetUserInfo(ctx context.Context) (UserView, error) {
	return sf.Get(ctx, jwtauth.FromUserId(ctx))
}

func (cUser) GetWithName(_ context.Context, name string) (item User, err error) {
	err = deployed.DB.Scopes(UserDB()).
		Where("username=? ", name).First(&item).Error
	return
}

func (cUser) GetWithDeptId(id int) (items []UserView, err error) {
	err = deployed.DB.Scopes(UserDB()).
		Select([]string{"sys_user.*", "sys_role.role_name"}).
		Joins("left join sys_role on sys_user.role_id=sys_role.role_id").
		Where("dept_id=?", id).
		Find(&items).Error
	return
}

func (cUser) QueryPage(ctx context.Context, qp UserQueryParam) ([]UserPage, paginator.Info, error) {
	var err error
	var items []UserPage

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
	db, err = GetDataScope("sys_user", jwtauth.FromUserId(ctx), db)
	if err != nil {
		return nil, paginator.Info{}, err
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

// 添加
func (cUser) Create(ctx context.Context, item User) (User, error) {
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

func (cUser) BatchDelete(id []int) error {
	return deployed.DB.Scopes(UserDB()).
		Where("user_id in (?)", id).Delete(&User{}).Error
}

// 修改
func (cUser) Update(ctx context.Context, id int, up User) (item User, err error) {
	if err = deployed.DB.Scopes(UserDB()).First(&item, id).Error; err != nil {
		return
	}

	up.Updator = jwtauth.FromUserIdStr(ctx)
	if up.RoleId == 0 {
		up.RoleId = item.RoleId
	}
	if up.Password != "" {
		up.Password, err = deployed.Verify.Hash(up.Password, "")
		if err != nil {
			return
		}
	}
	// 参数1:是要修改的数据
	// 参数2:是修改的数据
	err = deployed.DB.Table(up.TableName()).
		Model(&item).Updates(&up).Error
	return
}

func (cUser) UpdateAvatar(ctx context.Context, avatar string) error {
	id := jwtauth.FromUserId(ctx)
	return deployed.DB.Scopes(UserDB()).
		Where("user_id=?", id).
		Updates(map[string]interface{}{
			"avatar":  avatar,
			"updator": cast.ToString(id),
		}).Error
}

func (sf cUser) UpdatePassword(ctx context.Context, pwd UpdateUserPwd) error {
	item, err := sf.get(ctx, jwtauth.FromUserId(ctx))
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

func (cUser) get(_ context.Context, id int) (item UserView, err error) {
	err = deployed.DB.Scopes(UserDB()).
		Select([]string{"sys_user.*", "sys_role.role_name"}).
		Joins("left join sys_role on sys_user.role_id=sys_role.role_id").
		Where("user_id=?", id).First(&item).Error
	return
}
