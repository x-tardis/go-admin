package deployed

import (
	"github.com/spf13/viper"
)

type AliAccessKey struct {
	ID     string
	Secret string
}

func ViperAliAccessKey() AliAccessKey {
	c := viper.Sub("aliyun.accessKey")
	if c == nil {
		c = viper.New()
	}
	cf := AliAccessKey{
		c.GetString("accessKeyId"),
		c.GetString("accessKeySecret"),
	}
	return cf
}
