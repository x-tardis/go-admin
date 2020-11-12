package deployed

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/infra"
)

var FeatureConfig = new(Feature)
var AppConfig = new(Application)
var SslConfig = new(Ssl)
var GenConfig = new(Gen)

func init() {
	RegisterViperDefault(ViperApplicationDefault)
	RegisterViperDefault(ViperLimiterDefault)
	RegisterViperDefault(ViperCorsDefault)
	RegisterViperDefault(ViperJwtDefault)
	RegisterViperDefault(ViperRedisDefault)
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
	return AppConfig.Mode == infra.ModeDebug
}

func IsModeProd() bool {
	return AppConfig.Mode == infra.ModeProd
}

func IsModeDev() bool {
	return AppConfig.Mode == infra.ModeDev
}
