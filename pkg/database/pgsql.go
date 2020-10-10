package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/x-tardis/go-admin/common/config"
	"github.com/x-tardis/go-admin/common/global"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/textcolor"
	toolsConfig "github.com/x-tardis/go-admin/tools/config"
)

type PgSql struct{}

func (e *PgSql) Setup() {
	var err error

	global.Source = e.GetConnect()
	log.Println(global.Source)
	db, err := sql.Open("postgresql", global.Source)
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
		log.Fatalf("%s connect error %v", e.GetDriver(), err)
	} else {
		log.Printf("%s connect success!", e.GetDriver())
	}

	if deployed.DB.Error != nil {
		log.Fatalf("database error %v", deployed.DB.Error)
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
func (e *PgSql) Open(db *sql.DB, cfg *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{Conn: db}), cfg)
}

func (e *PgSql) GetConnect() string {
	return toolsConfig.DatabaseConfig.Source
}

func (e *PgSql) GetDriver() string {
	return toolsConfig.DatabaseConfig.Driver
}
