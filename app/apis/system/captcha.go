package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/servers"
)

func GenerateCaptchaHandler(c *gin.Context) {
	id, b64s, err := deployed.Captcha.Generate()
	if err != nil {
		servers.Fail(c, 500, "验证码获取失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}
