package system

import (
	"context"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

const (
	userIdKey    = "userId"
	usernameKey  = "username"
	roleIdKey    = "roleId"
	roleNameKey  = "roleName"
	roleKey      = "roleKey"
	dataScopeKey = "dataScope"
)

// NewJWTAuth new jwt auth
func NewJWTAuth(c *jwtauth.Config) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:      c.Realm,
		Key:        []byte(c.SecretKey),
		Timeout:    c.Timeout,
		MaxRefresh: c.MaxRefresh,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(jwtauth.Identities); ok {
				return jwt.MapClaims{
					userIdKey:    v.UserId,
					usernameKey:  v.Username,
					roleIdKey:    v.RoleId,
					roleNameKey:  v.RoleName,
					roleKey:      v.RoleKey,
					dataScopeKey: v.DataScope,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			identity := jwtauth.Identities{
				UserId:    cast.ToInt(claims[userIdKey]),
				Username:  cast.ToString(claims[usernameKey]),
				RoleId:    cast.ToInt(claims[roleIdKey]),
				RoleName:  cast.ToString(claims[roleNameKey]),
				RoleKey:   cast.ToString(claims[roleKey]),
				DataScope: cast.ToString(claims[dataScopeKey]),
			}
			ctx := context.WithValue(c.Request.Context(), jwtauth.IdentityKey{}, identity)
			c.Request = c.Request.WithContext(ctx)
			return identity
		},
		Authenticator: authenticator,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			_, ok := data.(jwtauth.Identities)
			return ok
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			servers.JSON(c, code, servers.Response{Code: code, Msg: message, Data: "{}"})
		},
		LogoutResponse: logoutResponse,
		TokenLookup:    "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:  "Bearer",
		TimeFunc:       time.Now,
	})
}

// @tags auth
// @summary 登陆
// @description 登陆,获取token
// @description 注意：开发模式：需要注意全部字段不能为空，账号,密码外可以传入0值
// @accept json
// @produce json
// @param account body models.Login  true "account"
// @success 200 {object} string "{"code": 200, "expire": "2019-08-07T12:45:48+08:00", "token": "xxxx" }"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /login [post]
func authenticator(c *gin.Context) (interface{}, error) {
	login := models.Login{}
	if err := c.ShouldBindJSON(&login); err != nil {
		loginLogRecord(c, false, "数据解析失败", login.Username)
		return nil, jwt.ErrMissingLoginValues
	}

	if deployed.AppConfig.Mode == infra.ModeProd &&
		!deployed.Captcha.Verify(login.CID, login.CCode, true) {
		loginLogRecord(c, false, "验证码错误", login.Username)
		return nil, jwt.ErrFailedAuthentication
	}
	user, role, err := login.GetUser()
	if err != nil {
		loginLogRecord(c, false, "登录失败", login.Username)
		deployed.RequestLogger.Debug(err.Error())
		return nil, jwt.ErrFailedAuthentication
	}
	if user.Status == models.StatusDisable {
		return nil, jwt.ErrForbidden
	}

	loginLogRecord(c, true, "登录成功", login.Username)
	return jwtauth.Identities{
		UserId:    user.UserId,
		Username:  user.Username,
		RoleId:    role.RoleId,
		RoleName:  role.RoleName,
		RoleKey:   role.RoleKey,
		DataScope: role.DataScope,
	}, nil
}

// @tags auth
// @summary 退出登录
// @description 退出登录
// @security Bearer
// @accept json
// @produce json
// @success 200 {object} string "{"code": 200, "msg": "成功退出系统" }"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /logout [post]
func logoutResponse(c *gin.Context, code int) {
	loginLogRecord(c, true, "退出成功", jwtauth.FromUserName(gcontext.Context(c)))
	servers.OK(c, servers.WithMsg("退出成功"))
}

// loginLogRecord Write log to database
func loginLogRecord(c *gin.Context, success bool, msg string, username string) {
	status := "0"
	if !success {
		status = "1"
	}
	if deployed.FeatureConfig.LoginDB.Load() {
		ua := user_agent.New(c.Request.UserAgent())
		browserName, browserVersion := ua.Browser()
		location := deployed.IPLocation(c.ClientIP())
		loginLog := models.LoginLog{
			Username:  username,
			Status:    status,
			Ip:        c.ClientIP(),
			Location:  location,
			Browser:   browserName + " " + browserVersion,
			Os:        ua.OS(),
			Platform:  ua.Platform(),
			LoginTime: time.Now(),
			Remark:    c.Request.UserAgent(),
			Msg:       msg,
		}

		models.CLoginLog.Create(context.Background(), loginLog) // nolint: errcheck
	}
}
