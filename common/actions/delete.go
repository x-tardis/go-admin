package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/common/dto"
	"github.com/x-tardis/go-admin/common/models"
	"github.com/x-tardis/go-admin/pkg/logger"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/tools"
)

// DeleteAction 通用删除动作
func DeleteAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := tools.GetOrm(c)
		if err != nil {
			logger.Error(err)
			return
		}

		msgID := tools.GenerateMsgIDFromContext(c)
		//删除操作
		req := control.Generate()
		err = req.Bind(c)
		if err != nil {
			logger.Errorf("MsgID[%s] Bind error: %s", msgID, err)
			servers.FailWithRequestID(c, http.StatusUnprocessableEntity, "参数验证失败")
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			servers.FailWithRequestID(c, http.StatusInternalServerError, "模型生成失败")
			return
		}

		object.SetUpdateBy(tools.GetUserIdUint(c))

		//数据权限检查
		p := GetPermissionFromContext(c)

		db = db.WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).Delete(object)
		if db.Error != nil {
			logger.Errorf("MsgID[%s] Delete error: %s", msgID, err)
			servers.FailWithRequestID(c, http.StatusInternalServerError, "删除失败")
			return
		}
		if db.RowsAffected == 0 {
			servers.FailWithRequestID(c, http.StatusForbidden, "无权删除该数据")
			return
		}
		servers.OKWithRequestID(c, object.GetId(), "删除成功")
		c.Next()
	}
}
