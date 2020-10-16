package system

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// @Summary 分页部门列表数据
// @Description 分页列表
// @Tags 部门
// @Param name query string false "name"
// @Param id query string false "id"
// @Param position query string false "position"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/depts [get]
// @Security Bearer
func GetDeptList(c *gin.Context) {
	var Dept models.SysDept

	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.DeptId = cast.ToInt(c.Request.FormValue("deptId"))
	Dept.DataScope = jwtauth.UserIdStr(c)
	result, err := Dept.SetDept(true)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func GetDeptTree(c *gin.Context) {
	var Dept models.SysDept
	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.DeptId = cast.ToInt(c.Request.FormValue("deptId"))
	result, err := Dept.SetDept(false)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 部门列表数据
// @Description 获取JSON
// @Tags 部门
// @Param deptId path string false "deptId"
// @Param position query string false "position"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/depts/{id} [get]
// @Security Bearer
func GetDept(c *gin.Context) {
	var Dept models.SysDept

	Dept.DeptId = cast.ToInt(c.Param("id"))
	Dept.DataScope = jwtauth.UserIdStr(c)
	result, err := Dept.Get()
	if err != nil {
		servers.Fail(c, 404, codes.NotFound)
		return
	}
	servers.OKWithRequestID(c, result, codes.GetSuccess)
}

// @Summary 添加部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param data body models.SysDept true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/depts [post]
// @Security Bearer
func InsertDept(c *gin.Context) {
	var data models.SysDept

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}
	data.CreateBy = jwtauth.UserIdStr(c)
	result, err := data.Create()
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
// @Param data body models.SysDept true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/depts [put]
// @Security Bearer
func UpdateDept(c *gin.Context) {
	var data models.SysDept

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	data.UpdateBy = jwtauth.UserIdStr(c)
	result, err := data.Update(data.DeptId)
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
func DeleteDept(c *gin.Context) {
	var data models.SysDept
	id, err := strconv.Atoi(c.Param("id"))
	_, err = data.Delete(id)
	if err != nil {
		servers.Fail(c, 500, "删除失败")
		return
	}
	servers.OKWithRequestID(c, "", codes.DeletedSuccess)
}

func GetDeptTreeRoleselect(c *gin.Context) {
	var Dept models.SysDept
	var SysRole models.SysRole
	id, err := strconv.Atoi(c.Param("roleId"))
	SysRole.RoleId = id
	result, err := Dept.SetDeptLabel()
	if err != nil {
		servers.Fail(c, -1, codes.NotFound)
		return
	}
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = SysRole.GetRoleDeptId()
		if err != nil {
			servers.Fail(c, -1, codes.NotFoundRelatedInfo)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":        200,
		"depts":       result,
		"checkedKeys": menuIds,
	})
}
