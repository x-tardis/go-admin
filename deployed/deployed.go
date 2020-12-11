package deployed

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v7"
	"github.com/mojocn/base64Captcha"
	"github.com/thinkgos/go-core-package/lib/password"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/middleware"
)

var Verify = new(password.BCrypt)

// var DriverString = base64Captcha.NewDriverString(46, 140, 2, 2, 5,
// 	"234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, []string{"wqy-microhei.ttc"}).ConvertFonts()
var Captcha = base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, base64Captcha.DefaultMemStore)

var CasbinEnforcer *casbin.SyncedEnforcer

var Redisc *redis.Client

var OSSClient *oss.Client
var OSSBucket *oss.Bucket

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

func SetupOSS() {
	var err error

	if FeatureConfig.OSS {
		OSSConfig = ViperAliyunOSS()
		OSSClient, err = oss.New(OSSConfig.Endpoint, OSSConfig.AccessKeyId, OSSConfig.AccessKeySecret)
		infra.HandlerError(err)
		OSSBucket, err = OSSClient.Bucket(OSSConfig.Bucket)
		infra.HandlerError(err)
	}
}
