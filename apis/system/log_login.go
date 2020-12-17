package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/go-core-package/lib/habit"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

// LoginLog api login log
type LoginLog struct{}

// @tags 登录日志
// @summary 获取登录日志列表
// @description 获取登录日志列表
// @security Bearer
// @accept json
// @produce json
// @param username query string false "username"
// @param ip query string false "ip"
// @param status query string false "status"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @Success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/loginlog [get]
func (LoginLog) QueryPage(c *gin.Context) {
	qp := models.LoginLogQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.CLoginLog.QueryPage(gcontext.Context(c), qp)
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

// @tags 登录日志
// @summary 获取登录日志
// @description 获取登录日志
// @security Bearer
// @accept json
// @produce json
// @param infoId path int true "infoId"
// @Success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/loginlog/{id} [get]
func (LoginLog) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CLoginLog.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 登录日志
// @summary 添加登录日志
// @description 添加登录日志
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.LoginLog true "newItem"
// @Success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/loginlog [post]
func (LoginLog) Create(c *gin.Context) {
	newItem := models.LoginLog{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CLoginLog.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.CreateFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 登录日志
// @summary 修改登录日志
// @description 修改登录日志
// @security Bearer
// @accept json
// @produce json
// @param up body models.LoginLog true "update item"
// @Success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/loginlog [put]
func (LoginLog) Update(c *gin.Context) {
	up := models.LoginLog{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	err := models.CLoginLog.Update(gcontext.Context(c), up.InfoId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.UpdateFailed),
			servers.WithError(err),
		)
		return
	}
	servers.OK(c, servers.WithMsg(prompt.UpdatedSuccess))
}

// @tags 登录日志
// @summary 批量删除登录日志
// @description 批量删除登录日志
// @security Bearer
// @accept json
// @produce json
// @param ids path string true "以逗号（,）分割的id列表,如果为clean,将清空日志"
// @Success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/loginlog/{ids} [delete]
func (LoginLog) BatchDelete(c *gin.Context) {
	var err error

	switch action := c.Param("ids"); action {
	case "clean":
		err = models.CLoginLog.Clean(gcontext.Context(c))
	default: // ids
		ids := habit.ParseIdsGroupInt(action)
		err = models.CLoginLog.BatchDelete(gcontext.Context(c), ids)
	}
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.DeleteFailed),
			servers.WithError(err),
		)
		return
	}
	servers.OK(c, servers.WithMsg(prompt.DeleteSuccess))
}
