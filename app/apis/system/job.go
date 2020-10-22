package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"

	"github.com/x-tardis/go-admin/app/service"
	dto2 "github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/servers"
)

type Job struct{}

// RemoveJobForService 调用service实现
func (Job) RemoveJobForService(c *gin.Context) {
	msgID := requestid.FromRequestID(c)
	var v dto2.GeneralDelDto
	err := c.BindUri(&v)
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg("参数验证失败"))
		return
	}
	s := service.SysJob{}
	s.MsgID = msgID
	s.Orm = dao.DB
	err = s.RemoveJob(&v)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithMsg(s.Msg))
}

// StartJobForService 启动job service实现
func (Job) StartJobForService(c *gin.Context) {
	msgID := requestid.FromRequestID(c)
	var v dto2.GeneralGetDto
	err := c.BindUri(&v)
	if err != nil {
		izap.Sugar.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg("参数验证失败"))
		return
	}
	s := service.SysJob{}
	s.Orm = dao.DB
	s.MsgID = msgID
	err = s.StartJob(&v)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithMsg(s.Msg))
}
