package config

import (
	"database/sql"
	"net/http"

	"go.uber.org/zap"

	"github.com/x-tardis/go-admin/pkg/izap"
)

type Config struct {
	saas   bool
	dbs    map[string]*DBConfig
	db     *DBConfig
	engine http.Handler
}

type DBConfig struct {
	Driver string
	DB     *sql.DB
}

// SetDbs 设置对应key的db
func (c *Config) SetDbs(key string, db *DBConfig) {
	c.dbs[key] = db
}

// GetDbs 获取所有map里的db数据
func (c *Config) GetDbs() map[string]*DBConfig {
	return c.dbs
}

// GetDbByKey 根据key获取db
func (c *Config) GetDbByKey(key string) *DBConfig {
	return c.dbs[key]
}

// SetDb 设置单个db
func (c *Config) SetDb(db *DBConfig) {
	c.db = db
}

// GetDb 获取单个db
func (c *Config) GetDb() *DBConfig {
	return c.db
}

// SetLogger 设置日志组件
func (c *Config) SetLogger(l *zap.SugaredLogger) {
	// logger.DefaultLogger = l
}

// GetLogger 获取日志组件
func (c *Config) GetLogger() *zap.SugaredLogger {
	return izap.Sugar
}

// SetSaas 设置是否是saas应用
func (c *Config) SetSaas(saas bool) {
	c.saas = saas
}

// GetSaas 获取是否是saas应用
func (c *Config) GetSaas() bool {
	return c.saas
}

func DefaultConfig() *Config {
	return &Config{}
}
