package deployed

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"

	"github.com/x-tardis/go-admin/pkg/database"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

var FeatureConfig = new(Feature)
var AppConfig = new(Application)
var DbConfig = new(database.Database)
var JwtConfig = new(jwtauth.Config)
var SslConfig = new(Ssl)
var GenConfig = new(Gen)
var CorsConfig = new(cors.Config)

// 载入配置文件
func SetupConfig(path string) {
	err := LoadConfig(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	AppConfig = ViperApplication()
	FeatureConfig = ViperFeature()
	DbConfig = ViperDatabase()
	JwtConfig = ViperJwt()
	SslConfig = ViperSsl()
	CorsConfig = ViperCors()
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
	viper.OnConfigChange(func(in fsnotify.Event) {
		// TODO: 防止重复操作
		c := viper.Sub("feature")
		FeatureConfig.DataScope.Store(c.GetBool("dataScope"))
		FeatureConfig.OperDB.Store(c.GetBool("operDB"))
		FeatureConfig.LoginDB.Store(c.GetBool("loginDB"))
	})
	viper.WatchConfig()
	return viper.ReadInConfig()
}
