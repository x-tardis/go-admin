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

// CreateAction 通用新增动作
func CreateAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := requestid.FromRequestID(c)
		//新增操作
		req := control.Generate()
		err := req.Bind(c)
		if err != nil {
			servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg("参数验证失败"))
			return
		}
		var object dto.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("模型生成失败"))
			return
		}
		object.SetCreator(uint(jwtauth.FromUserId(gcontext.Context(c))))
		err = dao.DB.WithContext(c).Create(object).Error
		if err != nil {
			izap.Sugar.Errorf("MsgID[%s] Create error: %s", msgID, err)
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("创建失败"))
			return
		}
		servers.OK(c, servers.WithData(object.GetId()), servers.WithMsg("创建成功"))
		c.Next()
	}
}
