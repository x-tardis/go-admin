package models

import (
	"github.com/thinkgos/go-core-package/extrand"
	"github.com/thinkgos/go-core-package/lib/password"
	"gorm.io/gorm"
)

var verify = new(password.BCrypt)

// BaseUser 密码登录基础用户
type BaseUser struct {
	Username     string `json:"username" gorm:"type:varchar(100);comment:用户名"`
	Salt         string `json:"-" gorm:"type:varchar(255);comment:加盐;<-"`
	PasswordHash string `json:"-" gorm:"type:varchar(128);comment:密码hash;<-"`
	Password     string `json:"password" gorm:"-"`
}

// SetPassword 设置密码
func (u *BaseUser) SetPassword(value string) {
	u.Password = value
	u.Salt = extrand.RandString(16)
	u.PasswordHash, _ = verify.Hash(u.Password, u.Salt)
}

// Verify 验证密码
func (u *BaseUser) Verify(db *gorm.DB, tableName string) bool {
	db.Table(tableName).Where("username = ?", u.Username).First(u)
	return verify.Compare(u.Password, u.Salt, u.PasswordHash) == nil
}
