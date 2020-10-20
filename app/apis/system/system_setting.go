package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type SysSetting struct{}

// @Summary 查询系统信息
// @Description 获取JSON
// @Tags 系统设置
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/system/setting [get]
func (SysSetting) Get(c *gin.Context) {
	item, err := models.CSetting.Get()
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @Tags 系统设置
// @Summary 更新设置
// @Description 更新设置
// @Param data body models.User true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/system/setting [post]
func (SysSetting) Update(c *gin.Context) {
	up := models.UpSetting{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CSetting.Update(models.Setting{
		Logo: up.Logo,
		Name: up.Name,
	})
	if err != nil {
		servers.Fail(c, http.StatusOK, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}
