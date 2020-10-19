package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/servers"
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
		servers.Fail(c, 500, codes.GetFail)
		return
	}
	servers.OKWithRequestID(c, item, codes.GetSuccess)
}

// @Tags 系统设置
// @Summary 更新设置
// @Description 更新设置
// @Param data body models.User true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/system/setting [post]
func (SysSetting) Update(c *gin.Context) {
	var up models.UpSetting

	if err := c.ShouldBindJSON(&up); err != nil {
		servers.FailWithRequestID(c, http.StatusOK, codes.DataParseFailed)
		return
	}

	a, err := models.CSetting.Update(models.Setting{
		Logo: up.Logo,
		Name: up.Name,
	})
	if err != nil {
		servers.FailWithRequestID(c, http.StatusOK, err.Error())
		return
	}
	servers.OKWithRequestID(c, a, codes.UpdatedSuccess)
}
