package migrate

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/models"
)

var version = make(map[int]func(db *gorm.DB, version string) error)
var mutex sync.Mutex

func init() {
	Register(Vs(0, 1, 0, 1), BaseModelsTables)
	Register(Vs(1, 1, 0, 1), BaseAdminData)
	Register(Vs(2, 1, 0, 1), BaseAdditionMenu)
	Register(Vs(3, 1, 0, 1), BaseFileDir)
}

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

// Vs 版本格式
//  1TTTMMMNNNFFF
// 1: 固定
// TT: 类型0-990, 主要用于migrate固定排序
// MMM,NNN,FFF: 主,次,修订版本0-999
func Vs(tp, major, minor, fixed int) int {
	return (((1000+tp)*1000+major)*1000+minor)*1000 + fixed
}
