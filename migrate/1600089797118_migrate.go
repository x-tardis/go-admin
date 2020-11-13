package migrate

import (
	"runtime"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	Register(GetFilename(fileName), _1600089797118Migrate)
}

func _1600089797118Migrate(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		f := &models.FileDir{
			Label:   "根目录",
			Pid:     0,
			Sort:    0,
			Path:    "",
			Creator: "1",
			Updator: "1",
		}
		err := tx.Create(f).Error
		if err != nil {
			return err
		}
		return tx.Create(&models.Migration{Version: version}).Error
	})
}
