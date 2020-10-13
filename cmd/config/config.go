package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

var configFile string
var StartCmd = &cobra.Command{
	Use:     "config",
	Short:   "Get Application config info",
	Example: "go-admin config -c config.yml",
	Run:     run,
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "Start server with provided configuration file")
}

func run(cmd *cobra.Command, args []string) {
	viper.BindPFlags(cmd.Flags()) // nolint: errcheck
	// viper.SetEnvPrefix("oam")
	// // OAM_CONFIGFILE
	// viper.BindEnv("config") // nolint: errcheck

	deployed.SetupConfig(configFile)

	application, errs := json.MarshalIndent(deployed.ApplicationConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("application:", string(application))

	jwt, errs := json.MarshalIndent(deployed.JwtConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("jwt:", string(jwt))

	database, errs := json.MarshalIndent(deployed.DatabaseConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("database:", string(database))

	gen, errs := json.MarshalIndent(deployed.GenConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("gen:", string(gen))

	loggerConfig, errs := json.MarshalIndent(deployed.LoggerConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("logger:", string(loggerConfig))
}
