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
// @description获取验证码
// @accept json
// @produce json
// @success 200 {object} servers.Response "成功"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/captcha [get]
func GetCaptcha(c *gin.Context) {
	id, b64s, err := deployed.Captcha.Generate()
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.CaptchaGetFailed),
			servers.WithError(err),
		)
		return
	}
	servers.OK(c, servers.WithData(gin.H{"id": id, "data": b64s}))
}
