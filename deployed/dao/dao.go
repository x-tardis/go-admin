package dao

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/x-tardis/go-admin/pkg/database"
)

var DB *gorm.DB

func SetupDatabase(c database.Config) {
	var err error

	DB, err = database.New(c, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		log.Fatalf("%s connect error %v", c.Dialect, err)
	}

	if c.LogMode {
		DB.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	}
}
