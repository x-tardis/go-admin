package migrate

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thinkgos/sharp/builder"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/migrate"
)

var configFile string
var Cmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Initialize the database",
	Example: fmt.Sprintf("%s migrate -c config/config.yaml", builder.Name),
	Run:     run,
}

func init() {
	Cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yaml", "Start server with provided configuration file")
}

func run(cmd *cobra.Command, _ []string) {
	log.Println(`start init...`)

	viper.BindPFlags(cmd.Flags()) // nolint: errcheck
	// viper.SetEnvPrefix("oam")
	// // OAM_CONFIGFILE
	// viper.BindEnv("config") // nolint: errcheck

	// 1. 读取配置
	deployed.SetupConfig(configFile)
	// 2. 设置日志
	deployed.SetupLogger()
	// 3. 初始化数据库链接
	dao.SetupDatabase(dao.DbConfig)
	// 4. 数据库迁移
	log.Println("数据库迁移开始")
	err := migrateModel()
	if err != nil {
		log.Fatalf("数据库基础数据初始化失败! %v ", err)
	}
	fmt.Println(`数据库基础数据初始化成功`)
}

func migrateModel() error {
	if dao.DbConfig.Driver == "mysql" {
		dao.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	return migrate.Migrate(dao.DB.Debug())
}
