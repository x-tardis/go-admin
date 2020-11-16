package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

// @tags 验证码
// @summary 获取验证码
// @description
// @accept json
// @produce json
// @success 200 {object} string "{"code": 200, "msg": "success", "data": "data", "id": "id"}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 404 {object} servers.Response "未找到相关信息"
// @failure 417 {object} servers.Response "客户端请求头错误"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/captcha [get]
func GetCaptcha(c *gin.Context) {
	id, b64s, err := deployed.Captcha.Generate()
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.CaptchaGetFailed))
		return
	}
	servers.JSON(c, http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  http.StatusText(http.StatusOK),
		"id":   id,
		"data": b64s,
	})
}
