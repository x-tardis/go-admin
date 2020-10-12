package system

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// @Summary 查询系统信息
// @Description 获取JSON
// @Tags 系统信息
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/setting [get]
func GetSetting(c *gin.Context) {
	var s models.SysSetting

	r, err := s.Get()
	if err != nil {
		servers.Fail(c, 500, codes.GetFail)
		return
	}
	if r.Logo != "" {
		if !strings.HasPrefix(r.Logo, "http") {
			r.Logo = fmt.Sprintf("http://%s/%s", c.Request.Host, r.Logo)
		}
	}
	servers.OKWithRequestID(c, r, codes.GetSuccess)
}

// @Summary 更新或提交系统信息
// @Description 获取JSON
// @Tags 系统信息
// @Param data body models.SysUser true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/system/setting [post]
func CreateSetting(c *gin.Context) {
	var s models.ResponseSystemConfig
	if err := c.ShouldBind(&s); err != nil {
		servers.FailWithRequestID(c, http.StatusOK, "缺少必要参数")
		return
	}

	var sModel models.SysSetting
	sModel.Logo = s.Logo
	sModel.Name = s.Name

	a, e := sModel.Update()
	if e != nil {
		servers.FailWithRequestID(c, http.StatusOK, e.Error())
		return
	}

	if a.Logo != "" {
		if !strings.HasPrefix(a.Logo, "http") {
			a.Logo = fmt.Sprintf("http://%s/%s", c.Request.Host, a.Logo)
		}
	}

	servers.OKWithRequestID(c, a, "提交成功")

}
