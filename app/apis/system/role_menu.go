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

// @Summary RoleMenu列表数据
// @Description 获取JSON
// @Tags 角色菜单
// @Param RoleId query string false "RoleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/rolemenu/{roleId} [get]
// @Security Bearer
func (RoleMenu) Get(c *gin.Context) {
	roleId := cast.ToInt(c.Param("roleId"))
	items, err := models.CRoleMenu.Get(gcontext.Context(c), roleId)
	if err != nil {
		servers.Fail(c, http.StatusOK, servers.WithPrompt(prompt.NotFound))
		return
	}
	servers.OK(c, servers.WithData(items))
}

func (RoleMenu) Create(c *gin.Context) {
	servers.OK(c, servers.WithMsg("添加成功"))
}

// @Summary 删除用户菜单数据
// @Description 删除数据
// @Tags 角色菜单
// @Param id path string true "id"
// @Param menu_id query string false "menu_id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/rolemenu/{roleid} [delete]
func DeleteRoleMenu(c *gin.Context) {
	roleId := cast.ToInt(c.Param("roleid"))
	err := models.CRoleMenu.DeleteWithRole(gcontext.Context(c), roleId)
	if err != nil {
		servers.Fail(c, http.StatusOK, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}
