package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/x-tardis/go-admin/cmd/migrate/migration"
	_ "github.com/x-tardis/go-admin/cmd/migrate/migration/version"
	"github.com/x-tardis/go-admin/common/models"
	"github.com/x-tardis/go-admin/pkg/deployed"
)

var configFile string
var StartCmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Initialize the database",
	Example: "go-admin migrate -c config.yaml",
	Run:     run,
}

// var exec bool

func init() {
	StartCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "Start server with provided configuration file")
	//StartCmd.PersistentFlags().BoolVarP(&exec, "exec", "e", false, "exec script")
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println(`start init`)

	viper.BindPFlags(cmd.Flags()) // nolint: errcheck
	// viper.SetEnvPrefix("oam")
	// // OAM_CONFIGFILE
	// viper.BindEnv("config") // nolint: errcheck

	//1. 读取配置
	deployed.SetupConfig(configFile)
	//2. 设置日志
	deployed.SetupLogger()
	//3. 初始化数据库链接
	deployed.SetupDatabase(deployed.DatabaseConfig.Driver, deployed.DatabaseConfig.Source)
	//4. 数据库迁移
	fmt.Println("数据库迁移开始")
	_ = migrateModel()
	//fmt.Println("数据库结构初始化成功！")
	//5. 数据初始化完成
	//if err := models.InitDb(); err != nil {
	//	global.Logger.Fatalf("数据库基础数据初始化失败！error: %v ", err)
	//}

	fmt.Println(`数据库基础数据初始化成功`)
}

func migrateModel() error {
	if deployed.DatabaseConfig.Driver == "mysql" {
		deployed.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	err := deployed.DB.Debug().AutoMigrate(&models.Migration{})
	if err != nil {
		return err
	}
	migration.Migrate.SetDb(deployed.DB.Debug())
	migration.Migrate.Migrate()
	return err
}
