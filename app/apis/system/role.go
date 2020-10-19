package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"

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
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	qp.Inspect()

	items, count, err := new(models.CallRole).QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, err.Error())
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
	call := new(models.CallRole)
	id := cast.ToInt(c.Param("id"))
	result, err := call.Get(gcontext.Context(c), id)
	menuIds := make([]int, 0)
	menuIds, err = call.GetMenuIds(id)
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
// @Router /api/v1/roles [post]
func (Role) Create(c *gin.Context) {
	newItem := models.SysRole{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, 500, codes.DataParseFailed)
		return
	}

	item, err := new(models.CallRole).Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	var t models.RoleMenu

	if len(item.MenuIds) > 0 {
		_, err = t.Insert(item.RoleId, newItem.MenuIds)
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
	servers.OKWithRequestID(c, item, codes.CreatedSuccess)
}

// @Summary 修改用户角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.SysRole true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/roles [put]
func (Role) Update(c *gin.Context) {
	up := models.SysRole{}

	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	item, err := new(models.CallRole).Update(gcontext.Context(c), up.RoleId, up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	var t models.RoleMenu
	_, err = t.DeleteRoleMenu(up.RoleId)
	if err != nil {
		servers.Fail(c, -1, "修改失败（delete rm）")
		return
	}
	if len(up.MenuIds) > 0 {
		_, err := t.Insert(up.RoleId, up.MenuIds)
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
	servers.OKWithRequestID(c, item, "修改成功")
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
	err := new(models.CallRole).BatchDelete(gcontext.Context(c), ids)
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

func UpdateRoleDataScope(c *gin.Context) {
	up := models.SysRole{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	item, err := new(models.CallRole).Update(gcontext.Context(c), up.RoleId, up)

	var t models.SysRoleDept
	_, err = t.DeleteRoleDept(up.RoleId)
	if err != nil {
		servers.Fail(c, -1, codes.CreatedFail)
		return
	}
	if up.DataScope == "2" {
		_, err := t.Insert(up.RoleId, up.DeptIds)
		if err != nil {
			servers.Fail(c, -1, codes.CreatedFail)
			return
		}
	}
	servers.OKWithRequestID(c, item, codes.UpdatedSuccess)
}
