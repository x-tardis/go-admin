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

// Post post
type Post struct{}

// @tags 岗位/Post
// @summary 获取岗位列表
// @description 获取岗位列表
// @security Bearer
// @accept  json
// @produce json
// @param postName query string false "postName"
// @param postCode query string false "postCode"
// @param status query string false "status"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @router /api/v1/posts [get]
func (Post) QueryPage(c *gin.Context) {
	qp := models.PostQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.CPost.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(&paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @tags 岗位/Post
// @summary 获取id岗位信息
// @description 获取id岗位信息
// @security Bearer
// @accept  json
// @produce json
// @param id path int true "post id"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @router /api/v1/posts/{id} [get]
func (Post) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CPost.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 岗位/Post
// @summary 添加岗位
// @description 获取JSON
// @security Bearer
// @accept  json
// @produce json
// @param newItem body models.Post true "new item"
// @success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @router /api/v1/posts [post]
func (Post) Create(c *gin.Context) {
	newItem := models.Post{}
	err := c.ShouldBindJSON(&newItem)
	if err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CPost.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 岗位/Post
// @summary 修改岗位
// @description 获取JSON
// @security Bearer
// @accept  json
// @produce json
// @param up body models.Post true "update item"
// @success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @router /api/v1/posts [put]
func (Post) Update(c *gin.Context) {
	up := models.Post{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	err := models.CPost.Update(gcontext.Context(c), up.PostId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.UpdateFailed))
		return
	}
	servers.OK(c)
}

// @tags 岗位/Post
// @summary 批量删除岗位
// @description 批量删除数据
// @security Bearer
// @accept  json
// @produce json
// @param ids path int true "id 列表,以','分隔"
// @success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @success 500 {string} string	"{"code": 500, "message": "删除失败"}"
// @router /api/v1/posts/{ids} [delete]
func (Post) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CPost.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithMsg(prompt.DeleteSuccess))
}
