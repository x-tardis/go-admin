package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/x-tardis/go-admin/common/config"
	"github.com/x-tardis/go-admin/common/global"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/textcolor"
	toolsConfig "github.com/x-tardis/go-admin/tools/config"
)

type Mysql struct{}

func (e *Mysql) Setup() {
	global.Source = e.GetConnect()
	global.Logger.Info(textcolor.Green(global.Source))
	db, err := sql.Open("mysql", global.Source)
	if err != nil {
		global.Logger.Fatal(textcolor.Red(e.GetDriver()+" connect error :"), err)
	}
	global.Cfg.SetDb(&config.DBConfig{
		Driver: "mysql",
		DB:     db,
	})
	deployed.DB, err = e.Open(db, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		global.Logger.Fatal(textcolor.Red(e.GetDriver()+" connect error :"), err)
	} else {
		global.Logger.Info(textcolor.Green(e.GetDriver() + " connect success !"))
	}

	if deployed.DB.Error != nil {
		global.Logger.Fatal(textcolor.Red(" database error :"), deployed.DB.Error)
	}

	if toolsConfig.LoggerConfig.EnabledDB {
		deployed.DB.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	}
}

// 打开数据库连接
func (e *Mysql) Open(db *sql.DB, cfg *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{Conn: db}), cfg)
}

// 获取数据库连接
func (e *Mysql) GetConnect() string {
	return toolsConfig.DatabaseConfig.Source
}

func (e *Mysql) GetDriver() string {
	return toolsConfig.DatabaseConfig.Driver
}
