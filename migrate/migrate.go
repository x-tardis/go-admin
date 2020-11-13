package migrate

import (
	"context"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/models"
)

var version = make(map[int]func(db *gorm.DB, version string) error)
var mutex sync.Mutex

func Register(key int, f func(db *gorm.DB, version string) error) {
	mutex.Lock()
	defer mutex.Unlock()
	if f == nil {
		panic("migrate: Register function is nil")
	}
	if _, dup := version[key]; dup {
		panic("migrate: Register called twice for key " + strconv.Itoa(key))
	}
	version[key] = f
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Migration{}); err != nil {
		return err
	}

	// clone
	verSlice := make([]int, 0, len(version))
	verMap := make(map[int]func(db *gorm.DB, version string) error)
	mutex.Lock()
	for k, v := range version {
		verMap[k] = v
		verSlice = append(verSlice, k)
	}
	mutex.Unlock()
	sort.Sort(sort.IntSlice(verSlice))

	for _, v := range verSlice {
		count, err := models.GetMigrationCount(context.Background(), db, v)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		err = verMap[v](db, strconv.Itoa(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFilename(s string) int {
	s = filepath.Base(s)
	return cast.ToInt(s[:13])
}
