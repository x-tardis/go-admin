package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/common/dto"
	"github.com/x-tardis/go-admin/common/models"
	"github.com/x-tardis/go-admin/logger"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/tools"
)

// CreateAction 通用新增动作
func CreateAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := tools.GetOrm(c)
		if err != nil {
			logger.Error(err)
			return
		}

		msgID := tools.GenerateMsgIDFromContext(c)
		//新增操作
		req := control.Generate()
		err = req.Bind(c)
		if err != nil {
			servers.FailWithRequestID(c, http.StatusUnprocessableEntity, "参数验证失败")
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			servers.FailWithRequestID(c, http.StatusInternalServerError, "模型生成失败")
			return
		}
		object.SetCreateBy(tools.GetUserIdUint(c))
		err = db.WithContext(c).Create(object).Error
		if err != nil {
			logger.Errorf("MsgID[%s] Create error: %s", msgID, err)
			servers.FailWithRequestID(c, http.StatusInternalServerError, "创建失败")
			return
		}
		servers.OKWithRequestID(c, object.GetId(), "创建成功")
		c.Next()
	}
}
