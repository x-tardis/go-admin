package log

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type LoginLog struct{}

// @Summary 登录日志列表
// @Description 获取JSON
// @Tags 登录日志
// @Param status query string false "status"
// @Param dictCode query string false "dictCode"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/loginlog [get]
// @Security Bearer
func (LoginLog) QueryPage(c *gin.Context) {
	qp := models.LoginLogQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	result, info, err := models.CLoginLog.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithError(err),
			servers.WithPrompt(prompt.QueryFailed))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(&paginator.Pages{
		Info: info,
		List: result,
	}))
}

// @Summary 通过编码获取登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Param infoId path int true "infoId"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/loginlog/{id} [get]
// @Security Bearer
func (LoginLog) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	result, err := models.CLoginLog.Get(id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
}

// @Summary 添加登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Accept  application/json
// @Product application/json
// @Param data body models.LoginLog true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/loginlog [post]
// @Security Bearer
func (LoginLog) Create(c *gin.Context) {
	newItem := models.LoginLog{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CLoginLog.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.CreateFailed),
			servers.WithError(err))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(item))
}

// @Summary 修改登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Accept  application/json
// @Product application/json
// @Param data body models.LoginLog true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/loginlog [put]
// @Security Bearer
func (LoginLog) Update(c *gin.Context) {
	up := models.LoginLog{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	result, err := models.CLoginLog.Update(gcontext.Context(c), up.InfoId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.UpdateFailed),
			servers.WithError(err))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
}

// @Summary 批量删除登录日志
// @Description 删除数据
// @Tags 登录日志
// @Param infoId path string true "以逗号（,）分割的infoId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/loginlog/{ids} [delete]
func (LoginLog) BatchDelete(c *gin.Context) {
	var err error

	action := c.Param("ids")
	switch action {
	case "clean":
		err = models.CLoginLog.Clean(gcontext.Context(c))
	default: // ids
		ids := infra.ParseIdsGroup(action)
		err = models.CLoginLog.BatchDelete(gcontext.Context(c), ids)
	}
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithPrompt(prompt.DeleteSuccess))
}
