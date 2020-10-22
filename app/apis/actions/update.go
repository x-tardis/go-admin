package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/app/models/dao"
	dto2 "github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// UpdateAction 通用更新动作
func UpdateAction(control dto2.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := requestid.FromRequestID(c)
		req := control.Generate()
		//更新操作
		err := req.Bind(c)
		if err != nil {
			servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg("参数验证失败"))
			return
		}
		var object dto2.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("模型生成失败"))
			return
		}
		object.SetUpdator(uint(jwtauth.FromUserId(gcontext.Context(c))))

		//数据权限检查
		p := GetPermissionFromContext(c)

		db := dao.DB.WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).Updates(object)
		if db.Error != nil {
			izap.Sugar.Errorf("MsgID[%s] Update error: %s", msgID, err)
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("更新失败"))
			return
		}
		if db.RowsAffected == 0 {
			servers.Fail(c, http.StatusForbidden, servers.WithMsg("无权更新该数据"))
			return
		}
		servers.OK(c, servers.WithData(object.GetId()), servers.WithMsg("更新成功"))
		c.Next()
	}
}
