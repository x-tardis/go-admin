package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/gin/gcontext"
	"github.com/thinkgos/x/lib/habit"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type OperLog struct{}

// @tags 操作日志
// @summary 获取操作日志列表
// @description 获取操作日志列表
// @security Bearer
// @accept json
// @produce json
// @param title query string false "title"
// @param operName query string false "operName"
// @param operIp query string false "operIp"
// @param businessType query string false "businessType"
// @param status query string false "status"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/operlog [get]
func (OperLog) QueryPage(c *gin.Context) {
	qp := models.OperLogQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.COperLog.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithError(err),
			servers.WithMsg(prompt.QueryFailed))
		return
	}
	servers.OK(c, servers.WithData(&paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @tags 操作日志
// @summary 通过id获取登录日志
// @description 通过id获取登录日志
// @security Bearer
// @accept json
// @produce json
// @Param id path int true "id"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/operlog/{id} [get]
func (OperLog) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.COperLog.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 操作日志
// @summary 添加操作日志
// @description 添加操作日志
// @security Bearer
// @accept json
// @produce json
// @Param newItem body models.OperLog true "new item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/operlog [post]
func (OperLog) Create(c *gin.Context) {
	newItem := models.OperLog{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	item, err := models.COperLog.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.CreateFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 操作日志
// @summary 批量删除操作日志
// @description 批量删除操作日志
// @security Bearer
// @accept json
// @produce json
// @param ids path string true "以逗号（,）分隔的dd"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/operlog/{ids} [delete]
func (OperLog) BatchDelete(c *gin.Context) {
	var err error

	action := c.Param("ids")
	switch action {
	case "clean":
		err = models.COperLog.Clean(gcontext.Context(c))
	default: // ids
		ids := habit.ParseIdsGroupInt(action)
		err = models.COperLog.BatchDelete(gcontext.Context(c), ids)
	}
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithMsg(prompt.DeleteSuccess))
}
