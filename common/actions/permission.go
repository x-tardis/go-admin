package actions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/middleware"
	"github.com/x-tardis/go-admin/pkg/servers"
)

type DataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
}

func PermissionAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := middleware.GetOrm(c)
		if err != nil {
			izap.Sugar.Error(err)
			return
		}

		msgID := middleware.GenerateMsgIDFromContext(c)
		var p = new(DataPermission)
		if userId := jwtauth.UserIdStr(c); userId != "" {
			p, err = newDataPermission(db, userId)
			if err != nil {
				izap.Sugar.Errorf("MsgID[%s] PermissionAction error: %s", msgID, err)
				servers.FailWithRequestID(c, http.StatusInternalServerError, "权限范围鉴定错误")
				c.Abort()
				return
			}
		}
		c.Set(PermissionKey, p)
		c.Next()
	}
}

func newDataPermission(tx *gorm.DB, userId interface{}) (*DataPermission, error) {
	var err error
	p := &DataPermission{}

	err = tx.Table("sys_user").
		Select("sys_user.user_id", "sys_role.role_id", "sys_user.dept_id", "sys_role.data_scope").
		Joins("left join sys_role on sys_role.role_id = sys_user.role_id").
		Where("sys_user.user_id = ?", userId).
		Scan(p).Error
	if err != nil {
		err = errors.New("获取用户数据出错 msg:" + err.Error())
		return nil, err
	}
	return p, nil
}

func Permission(tableName string, p *DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !deployed.AppConfig.EnableDP {
			return db
		}
		switch p.DataScope {
		case "2":
			return db.Where(tableName+".creator in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", p.RoleId)
		case "3":
			return db.Where(tableName+".creator in (SELECT user_id from sys_user where dept_id = ? )", p.DeptId)
		case "4":
			return db.Where(tableName+".creator in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+cast.ToString(p.DeptId)+"%")
		case "5":
			return db.Where(tableName+".creator = ?", p.UserId)
		default:
			return db
		}
	}
}

func getPermissionFromContext(c *gin.Context) *DataPermission {
	p := new(DataPermission)
	if pm, ok := c.Get(PermissionKey); ok {
		switch v := pm.(type) {
		case *DataPermission:
			p = v
		}
	}
	return p
}

// PermissionForNoAction 提供非action写法数据范围约束
func GetPermissionFromContext(c *gin.Context) *DataPermission {
	return getPermissionFromContext(c)
}
