package system

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

// @Summary RoleMenu列表数据
// @Description 获取JSON
// @Tags 角色菜单
// @Param RoleId query string false "RoleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/rolemenu [get]
// @Security Bearer
func GetRoleMenu(c *gin.Context) {
	var Rm models.RoleMenu
	err := c.ShouldBind(&Rm)
	result, err := Rm.Get()
	if err != nil {
		servers.Fail(c, http.StatusOK, servers.WithMsg("抱歉未找到相关信息"))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
}

type RoleMenuPost struct {
	RoleId   string
	RoleMenu []models.RoleMenu
}

func InsertRoleMenu(c *gin.Context) {
	servers.JSON(c, http.StatusOK, servers.WithMsg("添加成功"))
}

// @Summary 删除用户菜单数据
// @Description 删除数据
// @Tags 角色菜单
// @Param id path string true "id"
// @Param menu_id query string false "menu_id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/rolemenu/{id} [delete]
func DeleteRoleMenu(c *gin.Context) {
	var t models.RoleMenu
	id := c.Param("id")
	menuId := c.Request.FormValue("menu_id")
	fmt.Println(menuId)
	_, err := t.Delete(id, menuId)
	if err != nil {
		servers.Fail(c, http.StatusOK, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithPrompt(prompt.DeleteSuccess))
}
