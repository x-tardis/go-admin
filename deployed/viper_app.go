package deployed

import (
	"log"
	"net"
	"strings"

	"github.com/spf13/viper"
)

type Application struct {
	Mode          string // 工作模式
	Name          string // 应用名称
	Host          string // 主机名
	Port          string // 端口
	ReadTimeout   int    // 读超时
	WriterTimeout int    // 写超时
}

func (sf Application) Addr() string {
	return net.JoinHostPort(sf.Host, sf.Port)
}

func ViperApplication() *Application {
	return &Application{
		viper.GetString("mode"),
		viper.GetString("name"),
		viper.GetString("host"),
		viper.GetString("port"),
		viper.GetInt("readTimeout"),
		viper.GetInt("writerTimeout"),
	}
}

func CanonicalCDN() string {
	if CDNDomain == "" {
		log.Fatal("deployed.CanonicalCDN: CDN domain can not empty")
	}
	if !strings.HasPrefix(CDNDomain, "http://") &&
		!strings.HasPrefix(CDNDomain, "https://") {
		CDNDomain = "http://" + CDNDomain
	}
	if strings.HasSuffix(CDNDomain, "/") {
		return CDNDomain
	}
	return CDNDomain + "/"
}
