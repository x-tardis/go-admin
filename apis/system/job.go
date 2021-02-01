package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	"github.com/thinkgos/sharp/gin/gcontext"
	"go.uber.org/zap"

	"github.com/x-tardis/go-admin/app/service"
	"github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type Job struct{}

// RemoveJobForService 调用service实现
func (Job) RemoveJobForService(c *gin.Context) {
	msgID := requestid.FromRequestID(c)
	var v dto.GeneralDelDto
	err := c.BindUri(&v)
	if err != nil {
		zap.S().Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg(prompt.IncorrectRequestParam))
		return
	}
	s := service.SysJob{}
	err = s.RemoveJob(gcontext.Context(c), &v)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithMsg(s.Msg))
}

// StartJobForService 启动job service实现
func (Job) StartJobForService(c *gin.Context) {
	msgID := requestid.FromRequestID(c)
	var v dto.GeneralGetDto
	err := c.BindUri(&v)
	if err != nil {
		zap.S().Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg(prompt.IncorrectRequestParam))
		return
	}
	s := service.SysJob{}
	err = s.StartJob(gcontext.Context(c), &v)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithMsg(s.Msg))
}
