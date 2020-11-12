package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/codes"
)

type FileDir struct{}

// @tags filedir
// @summary 获取filedir树
// @description 获取filedir树
// @accept json
// @produce json
// @param id query int false "主键"
// @param label query string false "label名称"
// @param pId query int false "父id"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/filedir [get]
func (FileDir) QueryTree(c *gin.Context) {
	qp := models.FileDirQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	items, err := models.CFileDir.QueryTree(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(items))
}

// @tags filedir
// @summary 通过id获取filedir
// @description 通过id获取filedir
// @accept json
// @produce json
// @param id path int true "主键"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/filedir/{id} [get]
func (FileDir) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CFileDir.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound, servers.WithMsg(codes.NotFound))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags filedir
// @summary 创建filedir
// @description 创建filedir
// @accept json
// @produce json
// @param newItem body models.FileDir true "new item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/filedir [post]
func (FileDir) Create(c *gin.Context) {
	newItem := models.FileDir{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CFileDir.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags filedir
// @summary 更新filedir
// @description 更新filedir
// @accept json
// @produce json
// @param up body models.FileDir true "update item"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/filedir [put]
func (FileDir) Update(c *gin.Context) {
	up := models.FileDir{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusNotFound, servers.WithError(err))
		return
	}

	err := models.CFileDir.Update(gcontext.Context(c), up.Id, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(codes.UpdateFailed))
		return
	}
	servers.OK(c)
}

// @tags filedir
// @summary 批量删除filedir
// @description 批量删除filedir
// @accept json
// @produce json
// @param ids path string true "id列表,以','分隔"
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/filedir/{ids} [delete]
func (FileDir) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CFileDir.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(codes.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithMsg(codes.DeleteSuccess))
}
