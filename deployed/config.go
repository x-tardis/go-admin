package deployed

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/thinkgos/x/lib/habit"

	"github.com/x-tardis/go-admin/deployed/dao"
)

var FeatureConfig = new(Feature)
var AppConfig = new(Application)
var SslConfig = new(Ssl)
var GenConfig = new(Gen)

var AliAccessKeyConfig = AliAccessKey{}
var AliOSSConfig = AliyunOSS{}
var CDNDomain string // cdn域名,如果设置了oss,将使用oss的Bucket加速域名,否则使用配置的cdn

func init() {
	RegisterViperDefault(
		ViperLimiterDefault,
		ViperCorsDefault,
		ViperJwtDefault,
		ViperRedisDefault,
	)
}

// 载入配置文件
func SetupConfig(path string) {
	err := LoadConfig(path)
	if err != nil {
		log.Fatalf("Parse config file failed: %+v", err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		// TODO: 防止重复操作
		c := viper.Sub("feature")
		FeatureConfig.DataScope.Store(c.GetBool("dataScope"))
		FeatureConfig.OperDB.Store(c.GetBool("operDB"))
		FeatureConfig.LoginDB.Store(c.GetBool("loginDB"))
	})
	viper.WatchConfig()

	AppConfig = ViperApplication()
	FeatureConfig = ViperFeature()
	dao.DbConfig = dao.ViperDatabase()
	SslConfig = ViperSsl()
	GenConfig = ViperGen()
}

// 如果filename为空,将使用config.yaml配置文件,并在当前文件搜索
func LoadConfig(filename string) error {
	// 使用命令行或环境变量给的配置文件,否则使用默认路径进行查找
	if filename != "" {
		viper.SetConfigFile(filename)
	} else {
		viper.SetConfigName("config") // 文件名
		viper.SetConfigType("yaml")   // 配置类型
		viper.AddConfigPath(".")      // 增加搜索路径
	}

	ViperInitDefault()

	return viper.ReadInConfig()
}

func IsModeDebug() bool {
	return habit.IsModeDebug(AppConfig.Mode)
}

func IsModeProd() bool {
	return habit.IsModeProd(AppConfig.Mode)
}

func IsModeDev() bool {
	return habit.IsModeDev(AppConfig.Mode)
}
