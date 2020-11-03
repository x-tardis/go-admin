package database

import (
	"gorm.io/gorm"

	"gorm.io/gorm/schema"
)

// Config 数据库配置
type Config struct {
	Dialect  string `yaml:"dialect" json:"dialect"`
	UserName string `yaml:"userName" json:"userName"`
	Password string `yaml:"password" json:"password"`
	Protocol string `yaml:"protocol" json:"protocol"`
	Addr     string `yaml:"addr" json:"addr"`
	DbName   string `yaml:"dbName" json:"dbName"`
	LogMode  bool   `yaml:"logMode" json:"logMode"`
}

// Database 数据库配置
type Database struct {
	Driver    string
	Source    string
	EnableLog bool
}

func New(driver, source string) (*gorm.DB, error) {
	var dialect gorm.Dialector

	switch driver {
	case "mysql":
		dialect = newMsql(source)
	case "postgres":
		dialect = newPostgres(source)
	case "sqlite3":
		dialect = newSqlite3(source)
	default:
		panic("please select database driver one of [mysql|postgres|sqlite3], if use sqlite3, build tags with sqlite3!")
	}
	return gorm.Open(dialect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
}
