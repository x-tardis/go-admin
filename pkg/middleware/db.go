package middleware

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/x-tardis/go-admin/common/config"
)

func WithContextDb(dbMap map[string]*gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db, ok := dbMap["*"]; ok {
			c.Set("db", db)
		} else {
			c.Set("db", dbMap[c.Request.Host])
		}
		c.Next()
	}
}

func getGormFromDb(driver string, db *sql.DB, config *gorm.Config) (*gorm.DB, error) {
	switch driver {
	case "mysql":
		return gorm.Open(mysql.New(mysql.Config{Conn: db}), config)
	case "postgres":
		return gorm.Open(postgres.New(postgres.Config{Conn: db}), config)
	default:
		return nil, errors.New("not support this db driver")
	}
}

func GetGormFromConfig(cfg config.Conf) map[string]*gorm.DB {
	gormDB := make(map[string]*gorm.DB)
	if cfg.GetSaas() {
		var err error
		for k, v := range cfg.GetDbs() {
			gormDB[k], err = getGormFromDb(v.Driver, v.DB, &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})
			if err != nil {
				// deployed.Logger.Fatal(textcolor.Red(k+" connect error :"), err) TODO: mo
			}
		}
		return gormDB
	}
	c := cfg.GetDb()
	db, err := getGormFromDb(c.Driver, c.DB, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		// deployed.Logger.Fatal(textcolor.Red(c.Driver+" connect error :"), err) TODO: mo
	}
	gormDB["*"] = db
	return gormDB
}
