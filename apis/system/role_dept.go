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

type RoleDept struct{}

// @tags 部门/Dept
// @summary 获取部门树Label和角色已选的部门项
// @description
// @security Bearer
// @accept json
// @produce json
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/roleDeptTree/option/:roleId [get]
func (RoleDept) GetDeptTreeOptionRole(c *gin.Context) {
	roleId := cast.ToInt(c.Param("roleId"))
	tree, deptIds, err := models.CRoleDept.GetDeptTreeOption(gcontext.Context(c), roleId)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":        200,
		"depts":       tree,
		"checkedKeys": deptIds,
	})
}
