package deployed

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"

	"github.com/x-tardis/go-admin/pkg/database"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

var DbConfig = new(database.Database)
var AppConfig = new(Application)
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

	DbConfig = ViperDatabase()
	AppConfig = ViperApplication()
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

	return viper.ReadInConfig()
}
