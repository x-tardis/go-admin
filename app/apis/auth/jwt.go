package auth

import (
	"context"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/pkg/deployed"
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
			servers.OK(c, servers.WithCode(code), servers.WithMsg(message))
		},
		LogoutResponse: logoutResponse,
		TokenLookup:    "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:  "Bearer",
		TimeFunc:       time.Now,
	})
}

// @tags auth
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

	if err := c.ShouldBindJSON(&req); err != nil {
		loginLogRecord(c, false, "数据解析失败", req.Username)
		return nil, jwt.ErrMissingLoginValues
	}

	if deployed.AppConfig.Mode == infra.ModeProd &&
		!deployed.Captcha.Verify(req.UUID, req.Code, true) {
		loginLogRecord(c, false, "验证码错误", req.Username)
		return nil, jwt.ErrFailedAuthentication
	}
	user, role, err := req.GetUser()
	if err != nil {
		loginLogRecord(c, false, "登录失败", req.Username)
		deployed.RequestLogger.Debug(err.Error())
		return nil, jwt.ErrFailedAuthentication
	}
	loginLogRecord(c, true, "登录成功", req.Username)
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
	loginLogRecord(c, true, "退出成功", jwtauth.UserName(c))
	servers.OK(c, servers.WithMsg("退出成功"))
}

// loginLogRecord Write log to database
func loginLogRecord(c *gin.Context, success bool, msg string, username string) {
	status := "0"
	if !success {
		status = "1"
	}
	if deployed.EnabledDB {
		ua := user_agent.New(c.Request.UserAgent())
		browserName, browserVersion := ua.Browser()
		location := deployed.IPLocation(c.ClientIP())
		loginLog := models.LoginLog{
			Username:      username,
			Status:        status,
			Ipaddr:        c.ClientIP(),
			LoginLocation: location,
			Browser:       browserName + " " + browserVersion,
			Os:            ua.OS(),
			Platform:      ua.Platform(),
			LoginTime:     time.Now(),
			Remark:        c.Request.UserAgent(),
			Msg:           msg,
		}

		models.CLoginLog.Create(context.Background(), loginLog) // nolint: errcheck
	}
}
