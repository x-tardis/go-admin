package actions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/common/dto"
	"github.com/x-tardis/go-admin/common/models"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/middleware"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// ViewAction 通用详情动作
func ViewAction(control dto.Control, f func() interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := middleware.GetOrm(c)
		if err != nil {
			izap.Sugar.Error(err)
			return
		}

		msgID := requestid.FromRequestID(c)
		//查看详情
		req := control.Generate()
		err = req.Bind(c)
		if err != nil {
			servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg("参数验证失败"))
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("模型生成失败"))
			return
		}

		var rsp interface{}
		if f != nil {
			rsp = f()
		} else {
			rsp, _ = req.GenerateM()
		}

		//数据权限检查
		p := GetPermissionFromContext(c)

		err = db.Model(object).WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).First(rsp).Error

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			servers.Fail(c, http.StatusNotFound, servers.WithMsg("查看对象不存在或无权查看"))
			return
		}
		if err != nil {
			izap.Sugar.Errorf("MsgID[%s] View error: %s", msgID, err)
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("查看失败"))
			return
		}
		servers.JSON(c, http.StatusOK, servers.WithData(rsp), servers.WithMsg("查看成功"))
		c.Next()
	}
}
