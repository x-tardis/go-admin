package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// @tags tags
// @summary
// @description
// @accept json
// @produce json
// @success 200 {object} string "{"code": 200, "msg": "success", "data": "data", "id": "id"}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 404 {object} servers.Response "未找到相关信息"
// @failure 417 {object} servers.Response "客户端请求头错误"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router / [get]
func GetCaptcha(c *gin.Context) {
	id, b64s, err := deployed.Captcha.Generate()
	if err != nil {
		servers.Fail(c, 500, "验证码获取失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"id":   id,
		"data": b64s,
	})
}
