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

type RoleMenu struct{}

// @tags 菜单/Menu
// @Summary 获取角色对应的菜单id数组
// @Description 获取JSON
// @Param id path int true "id"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/roleMenu/menuids/{id} [get]
// @Security Bearer
func (RoleMenu) GetMenuIDS(c *gin.Context) {
	items, err := models.CRoleMenu.GetIdsWithRoleName(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(items))
}

func (RoleMenu) GetMenuTreeOptionRole(c *gin.Context) {
	roleId := cast.ToInt(c.Param("roleId"))
	items, menuIds, err := models.CRoleMenu.GetMenuTreeOption(gcontext.Context(c), roleId)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.JSON(c, http.StatusOK, gin.H{
		"code":        200,
		"menus":       items,
		"checkedKeys": menuIds,
	})
}
