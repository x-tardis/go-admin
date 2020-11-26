package deployed

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v7"
	"github.com/mojocn/base64Captcha"
	"github.com/thinkgos/go-core-package/lib/password"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/middleware"
)

var Verify = new(password.BCrypt)

// var DriverString = base64Captcha.NewDriverString(46, 140, 2, 2, 5,
// 	"234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, []string{"wqy-microhei.ttc"}).ConvertFonts()
var Captcha = base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, base64Captcha.DefaultMemStore)

var CasbinEnforcer *casbin.SyncedEnforcer

var Redisc *redis.Client

func SetupCasbin() {
	var err error

	CasbinEnforcer, err = middleware.NewCasbinSyncedEnforcerFromString(CasbinModeText, dao.DB)
	if err != nil {
		panic(err)
	}
}

func SetupRedis() {
	Redisc = redis.NewClient(ViperRedis(nil, nil))
}
