package sysjob

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/admin/service"
	"github.com/x-tardis/go-admin/common/dto"
	"github.com/x-tardis/go-admin/common/log"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/tools"
)

type SysJob struct{}

// RemoveJobForService 调用service实现
func (e *SysJob) RemoveJobForService(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Errorf("msgID[%s] error:%s", msgID, err)
		servers.FailWithRequestID(c, http.StatusInternalServerError, err.Error())
		return
	}
	var v dto.GeneralDelDto
	err = c.BindUri(&v)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		servers.FailWithRequestID(c, http.StatusUnprocessableEntity, "参数验证失败")
		return
	}
	s := service.SysJob{}
	s.MsgID = msgID
	s.Orm = db
	err = s.RemoveJob(&v)
	if err != nil {
		servers.FailWithRequestID(c, http.StatusInternalServerError, err.Error())
		return
	}
	servers.OKWithRequestID(c, nil, s.Msg)
}

// StartJobForService 启动job service实现
func (e *SysJob) StartJobForService(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Errorf("msgID[%s] error:%s", msgID, err)
		servers.FailWithRequestID(c, http.StatusInternalServerError, err.Error())
		return
	}
	var v dto.GeneralGetDto
	err = c.BindUri(&v)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		servers.FailWithRequestID(c, http.StatusUnprocessableEntity, "参数验证失败")
		return
	}
	s := service.SysJob{}
	s.Orm = db
	s.MsgID = msgID
	err = s.StartJob(&v)
	if err != nil {
		servers.FailWithRequestID(c, http.StatusInternalServerError, err.Error())
		return
	}
	servers.OKWithRequestID(c, nil, s.Msg)
}
