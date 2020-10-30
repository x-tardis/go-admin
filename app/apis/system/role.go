package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"

	"github.com/x-tardis/go-admin/models"
)

type Role struct{}

// @tags 角色/Role
// @summary 角色列表数据
// @description Get JSON
// @accept json
// @produce json
// @security Bearer
// @param roleName query string false "roleName"
// @param roleKey query string false "roleKey"
// @param status query string false "status"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @router /api/v1/roles [get]
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

	servers.OK(c, servers.WithData(&paginator.Pages{
		Info: count,
		List: items,
	}))
}

// @tags 角色/Role
// @summary 获取Role数据
// @description 获取Role数据
// @security Bearer
// @accept json
// @produce json
// @param id path string true "id"
// @success 200 {string} string "{"code": 200, "data": [...]}"
// @success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @router /api/v1/role/{id} [get]
func (Role) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CRole.Get(gcontext.Context(c), id)
	menuIds, err := models.CRole.GetMenuIds(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	item.MenuIds = menuIds

	servers.OK(c, servers.WithData(item))
}

// @tags 角色/Role
// @summary 创建角色
// @description 创建角色
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.Role true "new item"
// @success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @router /api/v1/roles [post]
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

	servers.OK(c, servers.WithData(item))
}

// @tags 角色/Role
// @Summary 修改角色信息
// @Description 修改角色信息
// @security Bearer
// @accept json
// @produce json
// @param up body models.Role true "update item"
// @success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @router /api/v1/roles [put]
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
	servers.OK(c, servers.WithData(item))
}

// @tags 角色/Role
// @summary 删除用户角色
// @description 删除用户角色
// @security Bearer
// @accept json
// @produce json
// @param ids path int true "id列表,以','分隔"
// @success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @router /api/v1/roles/{ids} [delete]
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
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}

func (Role) UpdateDataScope(c *gin.Context) {
	up := models.Role{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CRole.UpdateDataScope(gcontext.Context(c), up.RoleId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithData(item))
}
