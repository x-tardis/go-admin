package syscategory

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

type Category struct{}

// @tags 文章分类/Category
// @summary 获取分类列表
// @description 获取分类列表
// @security Bearer
// @accept json
// @produce json
// @param name query string false "name"
// @param status query string false "status"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/categories [get]
func (Category) QueryPage(c *gin.Context) {
	qp := models.CategoryQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.CCategory.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithError(err),
			servers.WithMsg(prompt.QueryFailed))
		return
	}
	servers.OK(c, servers.WithData(paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @tags 文章分类/Category
// @summary 获取分类
// @description 获取分类
// @security Bearer
// @accept json
// @produce json
// @param id path int true "主键"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/categories/:id [get]
func (Category) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CCategory.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 文章分类/Category
// @summary 添加分类
// @description 添加分类
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.Category true "new item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/categories [post]
func (Category) Create(c *gin.Context) {
	newItem := models.Category{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CCategory.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 文章分类/Category
// @summary 更新分类
// @description 更新分类
// @security Bearer
// @accept json
// @produce json
// @param up body models.Category true "update item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/categories [put]
func (Category) Update(c *gin.Context) {
	up := models.Category{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusNotFound, servers.WithError(err))
		return
	}

	err := models.CCategory.Update(gcontext.Context(c), up.Id, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.UpdateFailed))
		return
	}
	servers.OK(c)
}

// @tags 文章分类/Category
// @summary 批量删除分类
// @description 批量删除分类
// @security Bearer
// @accept json
// @produce json
// @param ids path string true "以','分隔的id列表"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/categories/{ids} [get]
func (Category) BatchDelete(c *gin.Context) {
	ids := habit.ParseIdsGroupInt(c.Param("ids"))
	err := models.CCategory.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithMsg(prompt.DeleteSuccess))
}
