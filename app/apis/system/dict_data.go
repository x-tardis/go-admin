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

type DictData struct{}

// @tags 字典数据/DictData
// @summary 获取字典数据列表
// @description 获取字典数据列表
// @security Bearer
// @accept json
// @produce json
// @param dictType query string false "dictType"
// @param dictLabel query string false "dictLabel"
// @param status query string false "status"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @success 200 {object} paginator.Pages "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/data [get]
func (DictData) QueryPage(c *gin.Context) {
	qp := models.DictDataQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	item, info, err := models.CDictData.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(&paginator.Pages{
		Info: info,
		List: item,
	}))
}

// @tags 字典数据/DictData
// @summary 通过id获取字典数据
// @description 通过id获取字典数据
// @security Bearer
// @accept json
// @produce json
// @param id path int true "字典数据主键"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/data/{id} [get]
func (DictData) Get(c *gin.Context) {
	dictId := cast.ToInt(c.Param("id"))
	item, err := models.CDictData.Get(gcontext.Context(c), dictId)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 字典数据/DictData
// @summary 添加字典数据
// @description 添加字典数据
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.DictData true "newItem"
// @success 200 {object} string	"{"code": 200, "msg": ""}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/data [post]
func (DictData) Create(c *gin.Context) {
	newItem := models.DictData{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CDictData.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 字典数据/DictData
// @summary 修改字典数据
// @description 修改字典数据
// @security Bearer
// @accept json
// @produce json
// @param up body models.DictData true "update item"
// @success 200 {object} string	"{"code": 200, "msg": ""}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/data [put]
func (DictData) Update(c *gin.Context) {
	up := models.DictData{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	err := models.CDictData.Update(gcontext.Context(c), up.DictId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
		return
	}
	servers.OK(c)
}

// @Summary 删除字典数据
// @Description 删除数据
// @Tags 字典数据/DictData
// @Param dictId path int true "dictId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dict/data/{dictIds} [delete]

// @tags 字典数据/DictData
// @summary 批量删除字典数据
// @description 批量删除字典数据
// @security Bearer
// @accept json
// @produce json
// @param ids path string true "id列表,以','分隔"
// @success 200 {string} string	"{"code": 200, "msg": ""}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/data/{ids} [delete]
func (DictData) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CDictData.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}

// @tags 字典数据/DictData
// @summary 通过字典类型获取字典数据
// @description 通过字典类型获取字典数据
// @security Bearer
// @accept json
// @produce json
// @param dictType path string true "dictType"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/dict/databyType/{dictType} [get]
func (DictData) GetWithType(c *gin.Context) {
	dictType := c.Param("dictType")
	result, err := models.CDictData.GetWithType(gcontext.Context(c), dictType)
	if err != nil {
		servers.Fail(c, http.StatusNotFound, servers.WithPrompt(prompt.NotFound))
		return
	}
	servers.OK(c, servers.WithData(result))
}
