package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

// Dept dept
type Dept struct{}

// @tags 部门/Dept
// @summary 分页部门列表数据
// @description 分页列表
// @security Bearer
// @accept json
// @produce json
// @param deptId query int false "deptId"
// @param deptName query string false "deptName"
// @param deptPath query string false "deptPath"
// @param Status query string false "Status"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/depts [get]
func (Dept) QueryPage(c *gin.Context) {
	qp := models.DeptQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	tree, err := models.CDept.QueryTree(gcontext.Context(c), qp, true)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(tree))
}

// @tags 部门/Dept
// @summary 分页部门列表数据
// @description 分页列表
// @security Bearer
// @accept json
// @produce json
// @param deptId query int false "deptId"
// @param deptName query string false "deptName"
// @param deptPath query string false "deptPath"
// @param Status query string false "Status"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/deptTree [get]
func (Dept) QueryTree(c *gin.Context) {
	qp := models.DeptQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	items, err := models.CDept.QueryTree(gcontext.Context(c), qp, false)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(items))
}

// @tags 部门/Dept
// @summary 获取指定Id信息
// @description 获取指定Id信息
// @security Bearer
// @accept json
// @produce json
// @param deptId path int false "deptId"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/depts/{id} [get]
func (Dept) Get(c *gin.Context) {
	deptId := cast.ToInt(c.Param("id"))
	item, err := models.CDept.Get(gcontext.Context(c), deptId)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 部门/Dept
// @summary 添加部门信息
// @description 添加部门信息
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.Dept true "new item"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/depts [post]
func (Dept) Create(c *gin.Context) {
	newItem := models.Dept{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	item, err := models.CDept.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 部门/Dept
// @summary 更新部门信息
// @description 更新部门信息
// @security Bearer
// @accept json
// @produce json
// @param id path int true "id"
// @param up body models.Dept true "up"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/depts [put]
func (Dept) Update(c *gin.Context) {
	up := models.Dept{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CDept.Update(gcontext.Context(c), up.DeptId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 部门/Dept
// @summary 删除部门信息
// @description 删除部门信息
// @security Bearer
// @accept json
// @produce json
// @param id path int true "id"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/depts [delete]
func (Dept) Delete(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	err := models.CDept.Delete(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}
