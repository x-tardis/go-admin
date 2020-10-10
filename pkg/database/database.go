package database

import (
	"log"
	"os"
	"time"

	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/tools/config"
	"gorm.io/gorm/logger"
)

func Setup(driver, source string) {
	var err error

	switch driver {
	case "mysql":
		deployed.DB, err = newMsql(source)
	case "postgres":
		deployed.DB, err = newPgSql(source)
	case "sqlite3":
		deployed.DB, err = newSqlite3(source)
	default:
		panic("please select database driver one of [mysql|postgres|mysqlite3], if use sqlite3,build tags with sqlite3!")
	}

	if err != nil {
		log.Fatalf("%s connect error %v", driver, err)
	} else {
		log.Printf("%s connect success!", driver)
	}

	if deployed.DB.Error != nil {
		log.Fatalf("database error %v", deployed.DB.Error)
	}

	if config.LoggerConfig.EnabledDB {
		deployed.DB.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	}
}
