package deployed

import (
	"log"
	"os"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/mojocn/base64Captcha"
	"github.com/robfig/cron/v3"
	"github.com/thinkgos/go-core-package/lib/password"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/go-sql-driver/mysql"

	"github.com/x-tardis/go-admin/pkg/database"
	"github.com/x-tardis/go-admin/pkg/middleware"
)

var Verify = new(password.BCrypt)

// var DriverString = base64Captcha.NewDriverString(46, 140, 2, 2, 5,
// 	"234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, []string{"wqy-microhei.ttc"}).ConvertFonts()
var Captcha = base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, base64Captcha.DefaultMemStore)

var Cron *cron.Cron

var DB *gorm.DB

var CasbinEnforcer *casbin.SyncedEnforcer

func SetupCasbin() {
	e, err := middleware.NewCasbinSyncedEnforcerFromString(CasbinModeText, DB)
	if err != nil {
		panic(err)
	}
	CasbinEnforcer = e
}

func SetupDatabase(driver, source string) {
	var err error

	DB, err = database.New(driver, source)
	if err != nil {
		log.Fatalf("%s connect error %v", driver, err)
	}

	if EnabledDB {
		DB.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	}
}
