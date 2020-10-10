package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

//载入配置文件
func Setup(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	cfgDatabase := viper.Sub("database")
	if cfgDatabase == nil {
		panic("No found database in the configuration")
	}
	DatabaseConfig = InitDatabase(cfgDatabase)

	cfgApplication := viper.Sub("application")
	if cfgApplication == nil {
		panic("No found application in the configuration")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	cfgJwt := viper.Sub("jwt")
	if cfgJwt == nil {
		panic("No found jwt in the configuration")
	}
	JwtConfig = InitJwt(cfgJwt)

	cfgLogger := viper.Sub("logger")
	if cfgLogger == nil {
		panic("No found logger in the configuration")
	}
	LoggerConfig = InitLog(cfgLogger)

	cfgSsl := viper.Sub("ssl")
	if cfgSsl == nil {
		// Ssl不是系统强制要求的配置，默认可以不用配置，将设置为关闭状态
		fmt.Println("warning config not found ssl in the configuration")
		SslConfig = new(Ssl)
		SslConfig.Enable = false
	} else {
		SslConfig = InitSsl(cfgSsl)
	}

	cfgGen := viper.Sub("gen")
	if cfgGen == nil {
		panic("No found gen")
	}
	GenConfig = InitGen(cfgGen)
}
