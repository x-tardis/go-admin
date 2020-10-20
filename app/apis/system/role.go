package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"

	"github.com/x-tardis/go-admin/app/models"
)

type Role struct{}

// @Tags 角色/Role
// @Summary 角色列表数据
// @Description Get JSON
// @Param roleName query string false "roleName"
// @Param roleKey query string false "roleKey"
// @Param status query string false "status"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/roles [get]
// @Security Bearer
func (Role) QueryPage(c *gin.Context) {
	qp := models.RoleQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, count, err := models.CRole.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}

	servers.JSON(c, http.StatusOK, servers.WithData(&paginator.Pages{
		Info: count,
		List: items,
	}))
}

// @Summary 获取Role数据
// @Description 获取JSON
// @Tags 角色/Role
// @Param roleId path string false "roleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/role/{id} [get]
// @Security Bearer
func (Role) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CRole.Get(gcontext.Context(c), id)
	menuIds := make([]int, 0)
	menuIds, err = models.CRole.GetMenuIds(id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	item.MenuIds = menuIds

	servers.JSON(c, http.StatusOK, servers.WithData(item))
}

// @Summary 创建角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.Role true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/roles [post]
func (Role) Create(c *gin.Context) {
	newItem := models.Role{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CRole.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}

	var t models.RoleMenu

	if len(item.MenuIds) > 0 {
		_, err = t.Insert(item.RoleId, newItem.MenuIds)
		if err != nil {
			servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
			return
		}
	}

	err = deployed.CasbinEnforcer.LoadPolicy()
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(item))
}

// @Summary 修改用户角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.Role true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/roles [put]
func (Role) Update(c *gin.Context) {
	up := models.Role{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CRole.Update(gcontext.Context(c), up.RoleId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
		return
	}

	var t models.RoleMenu
	_, err = t.DeleteRoleMenu(up.RoleId)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
		return
	}
	if len(up.MenuIds) > 0 {
		_, err := t.Insert(up.RoleId, up.MenuIds)
		if err != nil {
			servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
			return
		}
	}

	err = deployed.CasbinEnforcer.LoadPolicy()
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(item))
}

// @Summary 删除用户角色
// @Description 删除数据
// @Tags 角色/Role
// @Param roleId path int true "roleId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/roles/{ids} [delete]
func (Role) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CRole.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}

	err = deployed.CasbinEnforcer.LoadPolicy()
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithPrompt(prompt.DeleteSuccess))
}

func UpdateRoleDataScope(c *gin.Context) {
	up := models.Role{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CRole.Update(gcontext.Context(c), up.RoleId, up)

	err = models.CRoleDept.DeleteRoleDept(up.RoleId)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	if up.DataScope == "2" {
		err := models.CRoleDept.Create(up.RoleId, up.DeptIds)
		if err != nil {
			servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.CreateFailed))
			return
		}
	}
	servers.JSON(c, http.StatusOK, servers.WithData(item))
}
