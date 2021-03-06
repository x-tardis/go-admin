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

// UpdateAction 通用更新动作
func UpdateAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := requestid.FromRequestID(c)
		req := control.Generate()
		//更新操作
		err := req.Bind(c)
		if err != nil {
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

		db := dao.DB.WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).Updates(object)
		if db.Error != nil {
			zap.S().Errorf("MsgID[%s] Update error: %s", msgID, err)
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.UpdateFailed))
			return
		}
		if db.RowsAffected == 0 {
			servers.Fail(c, http.StatusForbidden, servers.WithMsg(prompt.PermissionDenied))
			return
		}
		servers.OK(c, servers.WithData(object.GetId()), servers.WithMsg(prompt.UpdatedSuccess))
		c.Next()
	}
}
