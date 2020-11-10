package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

// DictType api dict type
type DictType struct{}

// @tags 字典类型/DictType
// @summary 获取字典列表
// @description 获取字典列表
// @security Bearer
// @accept json
// @produce json
// @param dictName query string false "dictName"
// @param dictType query string false "dictType"
// @param status query string false "status"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @success 200 {object} paginator.Pages "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/type [get]
func (DictType) QueryPage(c *gin.Context) {
	qp := models.DictTypeQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.CDictType.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithError(err),
			servers.WithPrompt(prompt.QueryFailed))
		return
	}
	servers.OK(c, servers.WithData(paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @tags 字典类型/DictType
// @summary 通过字典id获取字典类型
// @description 通过字典id获取字典类型
// @security Bearer
// @accept json
// @produce json
// @param id path int true "字典类型编码"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/type/{id} [get]
func (DictType) Get(c *gin.Context) {
	dictId := cast.ToInt(c.Param("id"))
	item, err := models.CDictType.Get(gcontext.Context(c), dictId)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 字典类型/DictType
// @summary 创建字典类型
// @description 创建字典类型
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.DictType true "new item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/type [post]
func (DictType) Create(c *gin.Context) {
	newItem := models.DictType{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CDictType.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.CreateFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 字典类型/DictType
// @summary 修改字典类型
// @description 修改字典类型
// @security Bearer
// @accept json
// @produce json
// @param up body models.DictType true "update item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/type [put]
func (DictType) Update(c *gin.Context) {
	up := models.DictType{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	err := models.CDictType.Update(gcontext.Context(c), up.DictId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.UpdateFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c)
}

// @tags 字典类型/DictType
// @summary 批量删除字典类型
// @description 批量删除字典类型
// @security Bearer
// @accept json
// @produce json
// @param ids path string true "id列表,以','分隔"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/type/{ids} [delete]
func (DictType) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CDictType.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}

// @tags 字典类型/DictType
// @summary 获取字典类型列表数据
// @description 获取字典类型列表数据
// @security Bearer
// @accept json
// @produce json
// @param dictName query string false "dictName"
// @param dictType query string false "dictType"
// @param status query string false "status"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/typeoption [get]
func (DictType) GetOption(c *gin.Context) {
	qp := models.DictTypeQueryParam{}
	if err := c.ShouldBindQuery(qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	items, err := models.CDictType.Query(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(items))
}
