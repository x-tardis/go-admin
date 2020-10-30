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

func (RoleMenu) GetMenuTreeOptionRole(c *gin.Context) {
	roleId := cast.ToInt(c.Param("roleId"))
	items, menuIds, err := models.CRoleMenu.GetMenuTreeOption(gcontext.Context(c), roleId)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.JSON(c, http.StatusOK, gin.H{
		"code":        200,
		"menus":       items,
		"checkedKeys": menuIds,
	})
}
