package router

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

func OperLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" && deployed.EnabledDB {
			OperLogRecord(c, c.Writer.Status(), time.Since(startTime))
		}
	}
}

// 写入操作日志表
// 该方法后续即将弃用
func OperLogRecord(c *gin.Context, statusCode int, latencyTime time.Duration) {
	reqUri := c.Request.RequestURI
	reqMethod := c.Request.Method
	menuList, _ := models.CMenu.Query(gcontext.Context(c), models.MenuQueryParam{
		Path:   reqUri,
		Action: reqMethod,
	})
	clientIP := c.ClientIP()
	username := jwtauth.FromUserName(gcontext.Context(c))
	sysOperLog := models.OperLog{
		OperIp:        clientIP,
		OperLocation:  deployed.IPLocation(clientIP),
		Status:        cast.ToString(statusCode),
		OperName:      username,
		RequestMethod: reqMethod,
		OperUrl:       reqUri,
		Creator:       username,
		OperTime:      time.Now(),
		LatencyTime:   latencyTime.String(),
		UserAgent:     c.Request.UserAgent(),
	}

	if reqUri == "/login" {
		sysOperLog.BusinessType = "10"
		sysOperLog.Title = "用户登录"
		sysOperLog.OperName = "-"
	} else if strings.Contains(reqUri, "/api/v1/logout") {
		sysOperLog.BusinessType = "11"
	} else if strings.Contains(reqUri, "/api/v1/captcha") {
		sysOperLog.BusinessType = "12"
		sysOperLog.Title = "验证码"
	} else {
		if reqMethod == "POST" {
			sysOperLog.BusinessType = "1"
		} else if reqMethod == "PUT" {
			sysOperLog.BusinessType = "2"
		} else if reqMethod == "DELETE" {
			sysOperLog.BusinessType = "3"
		}
	}
	sysOperLog.Method = reqMethod
	if len(menuList) > 0 {
		sysOperLog.Title = menuList[0].Title
	}
	b, _ := c.Get("body")
	sysOperLog.OperParam, _ = jsoniter.MarshalToString(b)
	if c.Err() == nil {
		sysOperLog.Status = "0"
	} else {
		sysOperLog.Status = "1"
	}
	models.COperLog.Create(context.Background(), sysOperLog) // nolint: errcheck
}
