package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/servers"
)

type Dept struct{}

// @Tags 部门
// @Summary 分页部门列表数据
// @Description 分页列表
// @Param deptId query int false "deptId"
// @Param deptName query string false "deptName"
// @Param deptPath query string false "deptPath"
// @Param Status query string false "Status"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/depts [get]
// @Security Bearer
func (Dept) QueryPage(c *gin.Context) {
	qp := models.DeptQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	tree, err := models.CDept.QueryTree(gcontext.Context(c), qp, true)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, tree, "")
}

// @Tags 部门
// @Summary 分页部门列表数据
// @Description 分页列表
// @Param deptId query int false "deptId"
// @Param deptName query string false "deptName"
// @Param deptPath query string false "deptPath"
// @Param Status query string false "Status"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/deptTree [get]
// @Security Bearer
func (Dept) QueryTree(c *gin.Context) {
	qp := models.DeptQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	result, err := models.CDept.QueryTree(gcontext.Context(c), qp, false)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Tags 部门
// @Summary 部门列表数据
// @Description 获取JSON
// @Param deptId path string false "deptId"
// @Param position query string false "position"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/depts/{id} [get]
// @Security Bearer
func (Dept) Get(c *gin.Context) {
	deptId := cast.ToInt(c.Param("id"))
	item, err := models.CDept.Get(gcontext.Context(c), deptId)
	if err != nil {
		servers.Fail(c, 404, codes.NotFound)
		return
	}
	servers.OKWithRequestID(c, item, codes.GetSuccess)
}

// @Summary 添加部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param data body models.Dept true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/depts [post]
// @Security Bearer
func (Dept) Create(c *gin.Context) {
	newItem := models.Dept{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}
	result, err := models.CDept.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, codes.CreatedSuccess)
}

// @Summary 修改部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body models.Dept true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/depts [put]
// @Security Bearer
func (Dept) Update(c *gin.Context) {
	up := models.Dept{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	result, err := models.CDept.Update(gcontext.Context(c), up.DeptId, up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, codes.UpdatedSuccess)
}

// @Summary 删除部门
// @Description 删除数据
// @Tags 部门
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/depts/{id} [delete]
func (Dept) Delete(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	// TODO: bug 删除不到部门
	if err := models.CDept.Delete(gcontext.Context(c), id); err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, "", codes.DeletedSuccess)
}

func GetDeptTreeRoleselect(c *gin.Context) {
	result, err := models.CDept.QueryLabelTree(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, -1, codes.NotFound)
		return
	}
	roleId := cast.ToInt(c.Param("roleId"))
	menuIds := make([]int, 0)
	if roleId != 0 {
		menuIds, err = models.CRole.GetDeptIds(gcontext.Context(c), roleId)
		if err != nil {
			servers.Fail(c, -1, codes.NotFoundRelatedInfo)
			return
		}
	}
	servers.JSONs(c, http.StatusOK, gin.H{
		"code":        200,
		"depts":       result,
		"checkedKeys": menuIds,
	})
}
