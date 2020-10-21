package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"

	dto2 "github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// CreateAction 通用新增动作
func CreateAction(control dto2.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := requestid.FromRequestID(c)
		//新增操作
		req := control.Generate()
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
		object.SetCreator(uint(jwtauth.UserId(c)))
		err = deployed.DB.WithContext(c).Create(object).Error
		if err != nil {
			izap.Sugar.Errorf("MsgID[%s] Create error: %s", msgID, err)
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("创建失败"))
			return
		}
		servers.OK(c, servers.WithData(object.GetId()), servers.WithMsg("创建成功"))
		c.Next()
	}
}
