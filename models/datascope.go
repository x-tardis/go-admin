package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cast"
	"github.com/thinkgos/x/lib/textcolor"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/x-tardis/go-admin/deployed"
)

const (
	ScopeWhole     = "1" // 全部权限
	ScopeCustomize = "2" // 自定义权限
	ScopeDept      = "3" // 仅部门权限
	ScopeDeptBelow = "4" // 仅部门及以下权限
	ScopeMyself    = "5" // 仅本人数据权限
)

func DataScope(tabler schema.Tabler, userId int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !deployed.FeatureConfig.DataScope.Load() {
			fmt.Printf("%s\n", `数据权限已经为您`+textcolor.Green(`关闭`)+`，如需开启请参考配置文件字段说明`)
			return db
		}

		user, err := CUser.Get(context.Background(), userId)
		if err != nil {
			db.AddError(errors.New("获取用户数据出错 msg:" + err.Error())) // nolint: errcheck
			return db
		}

		role, err := CRole.Get(context.Background(), user.RoleId)
		if err != nil {
			db.AddError(errors.New("获取用户数据出错 msg:" + err.Error())) // nolint: errcheck
			return db
		}

		switch role.DataScope {
		case ScopeCustomize:
			db = db.Where(tabler.TableName()+".creator in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", user.RoleId)
		case ScopeDept:
			db = db.Where(tabler.TableName()+".creator in (SELECT user_id from sys_user where dept_id = ? )", user.DeptId)
		case ScopeDeptBelow:
			db = db.Where(tabler.TableName()+".creator in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+cast.ToString(user.DeptId)+"%")
		case ScopeMyself, "":
			db = db.Where(tabler.TableName()+".creator = ?", userId)
		}
		return db
	}
}
