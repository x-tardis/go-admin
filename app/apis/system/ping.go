package system

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/pkg/servers"
)

// @tags 系统信息
// @summary ping/pong test
// @description  ping/pong test
// @accept json
// @produce json
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 404 {object} servers.Response "未找到相关信息"
// @failure 417 {object} servers.Response "客户端请求头错误"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/system/ping [get]
func Ping(c *gin.Context) {
	servers.Success(c, servers.WithMessage("pong"))
}
