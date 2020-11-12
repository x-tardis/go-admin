package syscontent

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/codes"
)

type Content struct{}

// @tags 文章内容/Content
// @summary 获取内容列表
// @description 获取内容列表
// @security Bearer
// @accept json
// @produce json
// @param cateId query string false "cateId分类id"
// @param name query string false "name"
// @param status query string false "status"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/contents [get]
func (Content) QueryPage(c *gin.Context) {
	qp := models.ContentQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.CContent.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithError(err),
			servers.WithMsg(codes.QueryFailed))
		return
	}
	servers.OK(c, servers.WithData(paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @tags 文章内容/Content
// @summary 获取内容
// @description 获取内容
// @security Bearer
// @accept json
// @produce json
// @param id path int true "主键"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/contents/:id [get]
func (Content) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CContent.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(codes.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 文章内容/Content
// @summary 添加内容
// @description 添加内容
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.Content true "new item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/contents [post]
func (Content) Create(c *gin.Context) {
	newItem := models.Content{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CContent.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 文章内容/Content
// @summary 更新内容
// @description 更新内容
// @security Bearer
// @accept json
// @produce json
// @param up body models.Content true "update item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/contents [put]
func (Content) Update(c *gin.Context) {
	up := models.Content{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusNotFound, servers.WithError(err))
		return
	}

	err := models.CContent.Update(gcontext.Context(c), up.Id, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(codes.UpdateFailed))
		return
	}
	servers.OK(c)
}

// @tags 文章内容/Content
// @summary 批量删除内容
// @description 批量删除内容
// @security Bearer
// @accept json
// @produce json
// @param ids path string true "以','分隔的id列表"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/contents/{ids} [get]
func (Content) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CContent.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(codes.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithMsg(codes.DeleteSuccess))
}
