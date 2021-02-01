package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	"github.com/thinkgos/sharp/gin/gcontext"
	"go.uber.org/zap"

	"github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

// DeleteAction 通用删除动作
func DeleteAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := requestid.FromRequestID(c)
		//删除操作
		req := control.Generate()
		err := req.Bind(c)
		if err != nil {
			zap.S().Errorf("MsgID[%s] Bind error: %s", msgID, err)
			servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg(prompt.IncorrectRequestParam))
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
			zap.S().Errorf("MsgID[%s] Delete error: %s", msgID, err)
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.DeleteFailed))
			return
		}
		if db.RowsAffected == 0 {
			servers.Fail(c, http.StatusForbidden, servers.WithMsg(prompt.PermissionDenied))
			return
		}
		servers.OK(c, servers.WithData(object.GetId()), servers.WithMsg(prompt.DeleteSuccess))
		c.Next()
	}
}
