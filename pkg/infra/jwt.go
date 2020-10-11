package infra

import "time"

// JWTConfig jwt 配置信息
type JWTConfig struct {
	Realm      string        `yaml:"realm" json:"realm"`
	SecretKey  string        `yaml:"secretKey" json:"secretKey"`
	Timeout    time.Duration `yaml:"timeout" json:"timeout"`
	MaxRefresh time.Duration `yaml:"maxRefresh" json:"maxRefresh"`
}

// JWTIdentity jwt identity
type JWTIdentity struct {
	UserId    int
	UserName  string
	RoleId    int
	RoleName  string
	RoleKey   string
	DataScope string
}
