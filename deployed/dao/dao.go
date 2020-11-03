package dao

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/x-tardis/go-admin/pkg/database"
)

var DB *gorm.DB

func SetupDatabase(c *database.Database) {
	var err error

	DB, err = database.New(c.Driver, c.Source)
	if err != nil {
		log.Fatalf("%s connect error %v", c.Driver, err)
	}

	if c.EnableLog {
		DB.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	}
}
