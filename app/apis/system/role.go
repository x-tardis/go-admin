package system

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"

	"github.com/x-tardis/go-admin/app/models"
)

// @Summary 角色列表数据
// @Description Get JSON
// @Tags 角色/Role
// @Param roleName query string false "roleName"
// @Param status query string false "status"
// @Param roleKey query string false "roleKey"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/rolelist [get]
// @Security Bearer
func GetRoleList(c *gin.Context) {
	var data models.SysRole
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = strconv.Atoi(index)
	}

	data.RoleKey = c.Request.FormValue("roleKey")
	data.RoleName = c.Request.FormValue("roleName")
	data.Status = c.Request.FormValue("status")
	data.DataScope = jwtauth.UserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	servers.Success(c, servers.WithData(&paginator.Page{
		List:      result,
		Total:     count,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}))
}

// @Summary 获取Role数据
// @Description 获取JSON
// @Tags 角色/Role
// @Param roleId path string false "roleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/role [get]
// @Security Bearer
func GetRole(c *gin.Context) {
	var Role models.SysRole
	Role.RoleId = cast.ToInt(c.Param("roleId"))
	result, err := Role.Get()
	menuIds := make([]int, 0)
	menuIds, err = Role.GetRoleMeunId()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	result.MenuIds = menuIds
	servers.OKWithRequestID(c, result, "")

}

// @Summary 创建角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.SysRole true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/role [post]
func InsertRole(c *gin.Context) {
	var data models.SysRole
	data.CreateBy = jwtauth.UserIdStr(c)
	err := c.Bind(&data)
	if err != nil {
		servers.Fail(c, 500, codes.DataParseFailed)
		return
	}
	id, err := data.Insert()
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	data.RoleId = id

	var t models.RoleMenu
	if len(data.MenuIds) > 0 {
		_, err = t.Insert(id, data.MenuIds)
		if err != nil {
			servers.Fail(c, -1, err.Error())
			return
		}
	}

	err = deployed.CasbinEnforcer.LoadPolicy()
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	servers.OKWithRequestID(c, data, codes.CreatedSuccess)
}

// @Summary 修改用户角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.SysRole true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/role [put]
func UpdateRole(c *gin.Context) {
	var data models.SysRole
	data.UpdateBy = jwtauth.UserIdStr(c)
	err := c.Bind(&data)
	if err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	result, err := data.Update(data.RoleId)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	var t models.RoleMenu
	_, err = t.DeleteRoleMenu(data.RoleId)
	if err != nil {
		servers.Fail(c, -1, "修改失败（delete rm）")
		return
	}
	if len(data.MenuIds) > 0 {
		_, err := t.Insert(data.RoleId, data.MenuIds)
		if err != nil {
			servers.Fail(c, -1, "修改失败（insert）")
			return
		}
	}

	err = deployed.CasbinEnforcer.LoadPolicy()
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "修改成功")
}

func UpdateRoleDataScope(c *gin.Context) {
	var data models.SysRole
	data.UpdateBy = jwtauth.UserIdStr(c)
	err := c.Bind(&data)
	if err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	result, err := data.Update(data.RoleId)

	var t models.SysRoleDept
	_, err = t.DeleteRoleDept(data.RoleId)
	if err != nil {
		servers.Fail(c, -1, codes.CreatedFail)
		return
	}
	if data.DataScope == "2" {
		_, err := t.Insert(data.RoleId, data.DeptIds)
		if err != nil {
			servers.Fail(c, -1, codes.CreatedFail)
			return
		}
	}
	servers.OKWithRequestID(c, result, codes.UpdatedSuccess)
}

// @Summary 删除用户角色
// @Description 删除数据
// @Tags 角色/Role
// @Param roleId path int true "roleId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/role/{roleId} [delete]
func DeleteRole(c *gin.Context) {
	var Role models.SysRole
	Role.UpdateBy = jwtauth.UserIdStr(c)

	IDS := infra.ParseIdsGroup(c.Param("roleId"))
	_, err := Role.BatchDelete(IDS)
	if err != nil {
		servers.Fail(c, -1, codes.DeletedFail)
		return
	}

	err = deployed.CasbinEnforcer.LoadPolicy()
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	servers.OKWithRequestID(c, "", codes.DeletedSuccess)
}
