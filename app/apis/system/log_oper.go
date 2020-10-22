package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type OperLog struct{}

// @Tags 操作日志
// @Summary 登录日志列表
// @Description 获取JSON
// @Param status query string false "status"
// @Param dictId query string false "dictId"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/operlog [get]
// @Security Bearer
func (OperLog) QueryPage(c *gin.Context) {
	qp := models.OperLogQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	result, info, err := models.COperLog.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithError(err),
			servers.WithPrompt(prompt.QueryFailed))
		return
	}
	servers.OK(c, servers.WithData(&paginator.Pages{
		Info: info,
		List: result,
	}))
}

// @Tags 操作日志
// @Summary 通过编码获取登录日志
// @Description 获取JSON
// @Param infoId path int true "infoId"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/operlog/{id} [get]
// @Security Bearer
func (OperLog) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.COperLog.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @Tags 操作日志
// @Summary 添加操作日志
// @Description 获取JSON
// @Accept  application/json
// @Product application/json
// @Param data body models.OperLog true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/operlog [post]
// @Security Bearer
func (OperLog) Create(c *gin.Context) {
	newItem := models.OperLog{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	item, err := models.COperLog.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.CreateFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @Tags 操作日志
// @Summary 批量删除操作日志
// @Description 删除数据
// @Param operId path string true "以逗号（,）分割的operId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/operlog/{ids} [delete]
func (OperLog) BatchDelete(c *gin.Context) {
	var err error

	action := c.Param("ids")
	switch action {
	case "clean":
		err = models.COperLog.Clean(gcontext.Context(c))
	default: // ids
		ids := infra.ParseIdsGroup(action)
		err = models.COperLog.BatchDelete(gcontext.Context(c), ids)
	}
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}
