package auth

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/tools"
)

// @Summary 登陆
// @Description 获取token
// @Description LoginHandler can be used by clients to get a jwt token.
// @Description Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// @Description Reply will be of the form {"token": "TOKEN"}.
// @Description dev mode：It should be noted that all fields cannot be empty, and a value of 0 can be passed in addition to the account password
// @Description 注意：开发模式：需要注意全部字段不能为空，账号密码外可以传入0值
// @Accept  application/json
// @Product application/json
// @Param account body models.Login  true "account"
// @Success 200 {string} string "{"code": 200, "expire": "2019-08-07T12:45:48+08:00", "token": ".eyJleHAiOjE1NjUxNTMxNDgsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU2NTE0OTU0OH0.-zvzHvbg0A" }"
// @Router /login [post]
func authenticator(c *gin.Context) (interface{}, error) {
	var req models.Login

	if err := c.ShouldBind(&req); err != nil {
		LoginLogRecord(c, "1", "数据解析失败", req.Username)
		return nil, jwt.ErrMissingLoginValues
	}

	if deployed.ApplicationConfig.Mode != "dev" &&
		!deployed.Captcha.Verify(req.UUID, req.Code, true) {
		LoginLogRecord(c, "1", "验证码错误", req.Username)
		return nil, jwt.ErrFailedAuthentication
	}
	user, role, err := req.GetUser()
	if err != nil {
		LoginLogRecord(c, "1", "登录失败", req.Username)
		deployed.RequestLogger.Println(err.Error())
		return nil, jwt.ErrFailedAuthentication
	}
	LoginLogRecord(c, "0", "登录成功", req.Username)
	return infra.JWTIdentity{
		UserId:    user.UserId,
		UserName:  user.Username,
		RoleId:    role.RoleId,
		RoleName:  role.RoleName,
		RoleKey:   role.RoleKey,
		DataScope: role.DataScope,
	}, nil
}

// @Summary 退出登录
// @Description 获取token
// LoginHandler can be used by clients to get a jwt token.
// Reply will be of the form {"token": "TOKEN"}.
// @Accept  application/json
// @Product application/json
// @Success 200 {string} string "{"code": 200, "msg": "成功退出系统" }"
// @Router /logout [post]
// @Security Bearer
func logoutResponse(c *gin.Context, code int) {
	LoginLogRecord(c, "0", "退出成功", tools.GetUserName(c))
	c.JSON(http.StatusOK, gin.H{"code": code, "msg": "退出成功"})
}

// LoginLogRecord Write log to database
func LoginLogRecord(c *gin.Context, status string, msg string, username string) {
	if deployed.LoggerConfig.EnabledDB {
		ua := user_agent.New(c.Request.UserAgent())
		browserName, browserVersion := ua.Browser()
		location := tools.GetLocation(c.ClientIP())
		loginLog := models.LoginLog{
			Ipaddr:        c.ClientIP(),
			Username:      username,
			LoginLocation: location,
			LoginTime:     time.Now(),
			Status:        status,
			Remark:        c.Request.UserAgent(),
			Browser:       browserName + " " + browserVersion,
			Os:            ua.OS(),
			Msg:           msg,
			Platform:      ua.Platform(),
		}
		loginLog.Create() // nolint: errcheck
	}
}
