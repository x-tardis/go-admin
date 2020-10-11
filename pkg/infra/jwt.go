package infra

import "time"

const JwtPayloadKey = "JWT_PAYLOAD"

// Jwt 配置信息
type JWT struct {
	Realm      string        `yaml:"realm" json:"realm"`
	SecretKey  string        `yaml:"secretKey" json:"secretKey"`
	Timeout    time.Duration `yaml:"timeout" json:"timeout"`
	MaxRefresh time.Duration `yaml:"maxRefresh" json:"maxRefresh"`
}
