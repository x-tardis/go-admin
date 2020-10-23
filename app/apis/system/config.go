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

type Config struct{}

// @tags 系统配置
// @summary 获取系统配置
// @description 获取系统配置
// @security Bearer
// @accept json
// @produce json
// @param configKey query string false "系统配置的Key"
// @param configName query string false "系统配置的名称"
// @param configType query string false "系统配置是否内置,由用户指定值,服务只存储"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/configs [get]
func (Config) QueryPage(c *gin.Context) {
	qp := models.ConfigQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.CConfig.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(&paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @tags 系统配置
// @summary 获取配置
// @description 获取配置
// @security Bearer
// @accept json
// @produce json
// @Param id path int true "系统配置主键"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/configs/{id} [get]
func (Config) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CConfig.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 系统配置
// @summary 通过Key获取配置
// @description 通过Key获取配置
// @security Bearer
// @accept json
// @produce json
// @Param key path string true "系统配置的key"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/configKey/{key} [get]
func (Config) GetWithKey(c *gin.Context) {
	key := c.Param("key")
	item, err := models.CConfig.GetWithKey(gcontext.Context(c), key)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item), servers.WithMsg(item.ConfigValue))
}

// @tags 系统配置
// @summary 创建系统配置
// @description 创建系统配置
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.Config true "new item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/configs [post]
func (Config) Create(c *gin.Context) {
	newItem := models.Config{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CConfig.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 系统配置
// @summary 更新系统配置
// @description 更新系统配置
// @security Bearer
// @accept json
// @produce json
// @param up body models.Config true "update item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/configs [put]
func (Config) Update(c *gin.Context) {
	up := models.Config{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CConfig.Update(gcontext.Context(c), up.ConfigId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 系统配置
// @summary 批量删除系统设置
// @description 批量删除系统设置
// @security Bearer
// @accept json
// @produce json
// @Param ids path string true "id,以','分隔的id列表"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/configs/{ids} [delete]
func (Config) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CConfig.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}
