package deployed

import (
	"github.com/casbin/casbin/v2"
	"github.com/mojocn/base64Captcha"
	"github.com/robfig/cron/v3"
	"github.com/thinkgos/go-core-package/lib/password"

	_ "github.com/go-sql-driver/mysql"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/middleware"
)

var Verify = new(password.BCrypt)

// var DriverString = base64Captcha.NewDriverString(46, 140, 2, 2, 5,
// 	"234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, []string{"wqy-microhei.ttc"}).ConvertFonts()
var Captcha = base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, base64Captcha.DefaultMemStore)

var Cron *cron.Cron

var CasbinEnforcer *casbin.SyncedEnforcer

func SetupCasbin() {
	var err error

	CasbinEnforcer, err = middleware.NewCasbinSyncedEnforcerFromString(CasbinModeText, dao.DB)
	if err != nil {
		panic(err)
	}
}
