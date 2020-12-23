package database

import (
	"fmt"
	"log"

	"github.com/thinkgos/go-core-package/lib/univ"
	"gorm.io/gorm"

	"gorm.io/gorm/schema"
)

// Config 数据库配置
type Config struct {
	Dialect  string            `yaml:"dialect" json:"dialect"` // mysql sqlite3 postgres
	Username string            `yaml:"username" json:"username"`
	Password string            `yaml:"password" json:"password"`
	Protocol string            `yaml:"protocol" json:"protocol"`
	Host     string            `yaml:"host" json:"host"`
	Port     string            `yaml:"port" json:"port"`
	DbName   string            `yaml:"dbName" json:"dbName"`
	Extend   map[string]string `yaml:"extend" json:"extend"`
	LogMode  bool              `yaml:"logMode" json:"logMode"`
}

func New(c Config) (*gorm.DB, error) {
	var dialect gorm.Dialector

	switch c.Dialect {
	case "mysql":
		values := make(univ.Values)
		values.Add("charset", "utf8mb4")
		// values.Add("parseTime", "True")
		values.Add("loc", "Local")

		for k, v := range c.Extend {
			values.Add(k, v)
		}

		dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?%s",
			c.Username, c.Password, c.Protocol, c.Host, c.Port, c.DbName, values.Encode("=", "&")) // DSN data source name
		log.Println(dsn)
		dialect = newMsql(dsn)
	case "postgres":
		values := make(univ.Values)
		values.Add("user", c.Username)
		values.Add("password", c.Password)
		values.Add("host", c.Host)
		values.Add("port", c.Port)
		values.Add("dbname", c.DbName)
		for k, v := range c.Extend {
			values.Add(k, v)
		}
		dialect = newPostgres(values.Encode("=", " "))
	case "sqlite3":
		dialect = newSqlite3(c.DbName)
	default:
		panic("please select database driver one of [mysql|postgres|sqlite3], if use sqlite3, build tags with sqlite3!")
	}
	return gorm.Open(dialect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
}
