package migrate

import (
	"runtime"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	Register(GetFilename(fileName), _1599190683659Tables)
}

func _1599190683659Tables(db *gorm.DB, version string) error {
	err := db.AutoMigrate(baseModels()...)
	if err != nil {
		return err
	}
	return db.Create(&models.Migration{Version: version}).Error
}
