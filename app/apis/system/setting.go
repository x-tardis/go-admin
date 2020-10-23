package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type Setting struct{}

// @tags 系统设置
// @summary 查询系统设置
// @description 查询系统设置
// @accept json
// @produce json
// @success 200 {object} models.Setting "系统设置信息"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/system/setting [get]
func (Setting) Get(c *gin.Context) {
	item, err := models.CSetting.Get(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 系统设置
// @summary 更新系统设置
// @description 更新系统设置
// @security Bearer
// @accept json
// @produce json
// @param up body models.UpSetting true "更新的信息"
// @success 200 {object} models.Setting "系统设置信息"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/system/setting [post]
func (Setting) Update(c *gin.Context) {
	up := models.UpSetting{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CSetting.Update(gcontext.Context(c), up)
	if err != nil {
		servers.Fail(c, http.StatusOK, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}
