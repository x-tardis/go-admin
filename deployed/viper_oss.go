package deployed

import (
	"github.com/spf13/viper"
)

type AliyunOSS struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	BucketDomain    string
}

func ViperAliyunOSS() AliyunOSS {
	c := viper.Sub("oss")
	if c == nil {
		c = viper.New()
	}
	cf := AliyunOSS{
		c.GetString("endpoint"),
		c.GetString("accessKeyId"),
		c.GetString("accessKeySecret"),
		c.GetString("bucket"),
		"",
	}
	cf.BucketDomain = cf.Bucket + "." + cf.Endpoint
	return cf
}
