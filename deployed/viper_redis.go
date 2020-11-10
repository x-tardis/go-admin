package deployed

import (
	"crypto/tls"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

func ViperRedisDefault() {
	viper.SetDefault("redis.addr", "127.0.0.1:6379")
}

func ViperRedis(onConnect func(*redis.Conn) error, tlsConfig *tls.Config) *redis.Options {
	var tc *tls.Config

	c := viper.Sub("redis")
	if c.GetBool("enableTLS") {
		tc = tlsConfig
	}
	return &redis.Options{
		Network:            c.GetString("network"),
		Addr:               c.GetString("addr"),
		Password:           c.GetString("password"),
		DB:                 c.GetInt("db"),
		MaxRetries:         c.GetInt("maxRetries"),
		MinRetryBackoff:    c.GetDuration("minRetryBackoff"),
		MaxRetryBackoff:    c.GetDuration("maxRetryBackoff"),
		DialTimeout:        c.GetDuration("dialTimeout"),
		ReadTimeout:        c.GetDuration("readTimeout"),
		WriteTimeout:       c.GetDuration("writeTimeout"),
		PoolSize:           c.GetInt("poolSize"),
		MinIdleConns:       c.GetInt("minIdleConns"),
		MaxConnAge:         c.GetDuration("maxConnAge"),
		PoolTimeout:        c.GetDuration("poolTimeout"),
		IdleTimeout:        c.GetDuration("idleTimeout"),
		IdleCheckFrequency: c.GetDuration("idleCheckFrequency"),
		TLSConfig:          tc,
		OnConnect:          onConnect,
	}
}
