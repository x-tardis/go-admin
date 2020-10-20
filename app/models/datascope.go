package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cast"
	"github.com/thinkgos/go-core-package/lib/textcolor"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

func GetDataScope(tbname string, userId int, db *gorm.DB) (*gorm.DB, error) {
	if !deployed.AppConfig.EnableDP {
		fmt.Printf("%s\n", `数据权限已经为您`+textcolor.Green(`关闭`)+`，如需开启请参考配置文件字段说明`)
		return db, nil
	}

	user, err := CUser.Get(context.Background(), userId)
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}

	role, err := CRole.Get(context.Background(), user.RoleId)
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	if role.DataScope == "2" {
		db = db.Where(tbname+".creator in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", user.RoleId)
	}
	if role.DataScope == "3" {
		db = db.Where(tbname+".creator in (SELECT user_id from sys_user where dept_id = ? )", user.DeptId)
	}
	if role.DataScope == "4" {
		db = db.Where(tbname+".creator in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+cast.ToString(user.DeptId)+"%")
	}
	if role.DataScope == "5" || role.DataScope == "" {
		db = db.Where(tbname+".creator = ?", userId)
	}

	return db, nil
}

func DataScopes(tableName string, userId int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		user, err := CUser.Get(context.Background(), userId)
		if err != nil {
			db.Error = errors.New("获取用户数据出错 msg:" + err.Error())
			return db
		}
		role, err := CRole.Get(context.Background(), user.RoleId)
		if err != nil {
			db.Error = errors.New("获取用户数据出错 msg:" + err.Error())
			return db
		}
		if role.DataScope == "2" {
			return db.Where(tableName+".creator in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", user.RoleId)
		}
		if role.DataScope == "3" {
			return db.Where(tableName+".creator in (SELECT user_id from sys_user where dept_id = ? )", user.DeptId)
		}
		if role.DataScope == "4" {
			return db.Where(tableName+".creator in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+cast.ToString(user.DeptId)+"%")
		}
		if role.DataScope == "5" || role.DataScope == "" {
			return db.Where(tableName+".creator = ?", userId)
		}
		return db
	}
}
