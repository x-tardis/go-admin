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

type FileInfo struct{}

// @tags fileinfo
// @summary 获取fileinfo列表
// @description 获取fileinfo列表
// @accept json
// @produce json
// @param pId query int false "父id"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/fileinfo [get]
func (FileInfo) QueryPage(c *gin.Context) {
	qp := models.FileInfoQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.CFileInfo.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @tags fileinfo
// @summary 通过id获取fileinfo
// @description 通过id获取fileinfo
// @accept json
// @produce json
// @param id path int true "主键"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/fileinfo/{id} [get]
func (FileInfo) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CFileInfo.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound, servers.WithMsg(prompt.NotFound))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags fileinfo
// @summary 创建fileinfo
// @description 创建fileinfo
// @accept json
// @produce json
// @param newItem body models.FileInfo true "new item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/fileinfo [post]
func (FileInfo) Create(c *gin.Context) {
	newItem := models.FileInfo{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CFileInfo.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags fileinfo
// @summary 更新fileinfo
// @description 更新fileinfo
// @accept json
// @produce json
// @param up body models.FileInfo true "update item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/fileinfo [put]
func (FileInfo) Update(c *gin.Context) {
	up := models.FileInfo{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	err := models.CFileInfo.Update(gcontext.Context(c), up.Id, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.UpdateFailed))
		return
	}
	servers.OK(c)
}

// @tags fileinfo
// @summary 批量删除fileinfo
// @description 批量删除fileinfo
// @accept json
// @produce json
// @param ids path string true "id列表,以','分隔"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/fileinfo/{ids} [delete]
func (FileInfo) BatchDelete(c *gin.Context) {
	ids := habit.ParseIdsGroupInt(c.Param("ids"))
	err := models.CFileInfo.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithMsg(prompt.DeleteSuccess))
}
