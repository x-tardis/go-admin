package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// DeleteAction 通用删除动作
func DeleteAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := requestid.FromRequestID(c)
		//删除操作
		req := control.Generate()
		err := req.Bind(c)
		if err != nil {
			izap.Sugar.Errorf("MsgID[%s] Bind error: %s", msgID, err)
			servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg("参数验证失败"))
			return
		}
		var object dto.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("模型生成失败"))
			return
		}

		object.SetUpdator(uint(jwtauth.FromUserId(gcontext.Context(c))))

		//数据权限检查
		p := GetPermissionFromContext(c)

		db := dao.DB.WithContext(c).Scopes(Permission(object.TableName(), p)).
			Where(req.GetId()).Delete(object)
		if db.Error != nil {
			izap.Sugar.Errorf("MsgID[%s] Delete error: %s", msgID, err)
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("删除失败"))
			return
		}
		if db.RowsAffected == 0 {
			servers.Fail(c, http.StatusForbidden, servers.WithMsg("无权删除该数据"))
			return
		}
		servers.OK(c, servers.WithData(object.GetId()), servers.WithMsg("删除成功"))
		c.Next()
	}
}
