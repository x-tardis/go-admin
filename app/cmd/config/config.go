package config

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thinkgos/sharp/builder"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/deployed/dao"
)

var configFile string
var Cmd = &cobra.Command{
	Use:     "config",
	Short:   "Get Application config info",
	Example: fmt.Sprintf("%s config -c config/config.yml", builder.Name),
	Run:     run,
}

func init() {
	Cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yaml", "Start server with provided configuration file")
}

func run(cmd *cobra.Command, args []string) {
	viper.BindPFlags(cmd.Flags()) // nolint: errcheck
	// viper.SetEnvPrefix("oam")
	// // OAM_CONFIGFILE
	// viper.BindEnv("config") // nolint: errcheck

	deployed.SetupConfig(configFile)

	application, err := marshalIndentToString(deployed.AppConfig) // 转换成JSON返回的是byte[]
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("application:", application)

	jwt, err := marshalIndentToString(deployed.ViperJwt())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("jwt:", jwt)

	// cors, err := marshalIndentToString(deployed.CorsConfig)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("cors:", cors)

	database, err := marshalIndentToString(dao.DbConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("database:", database)

	gen, err := marshalIndentToString(deployed.GenConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("gen:", gen)

	loggerConfig, err := marshalIndentToString(deployed.ViperLogger())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("logger:", loggerConfig)
}

func marshalIndentToString(v interface{}) (string, error) {
	b, err := jsoniter.MarshalIndent(v, "", "   ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
