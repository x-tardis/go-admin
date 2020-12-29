package deployed

import (
	"github.com/spf13/viper"
)

type AliyunOSS struct {
	Endpoint string
	Bucket   string
	Https    bool
}

func ViperAliOSS() AliyunOSS {
	c := viper.Sub("aliyun.oss")
	if c == nil {
		c = viper.New()
	}
	cf := AliyunOSS{
		c.GetString("endpoint"),
		c.GetString("bucket"),
		c.GetBool("https"),
	}
	return cf
}
