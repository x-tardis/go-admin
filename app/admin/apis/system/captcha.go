package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/tools"
)

func GenerateCaptchaHandler(c *gin.Context) {
	id, b64s, err := deployed.Captcha.Generate()
	tools.HasError(err, "验证码获取失败", 500)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}
