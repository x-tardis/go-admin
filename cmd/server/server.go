package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thinkgos/go-core-package/lib/textcolor"
	"github.com/thinkgos/sharp/builder"

	"github.com/x-tardis/go-admin/app/jobs"
	"github.com/x-tardis/go-admin/app/router"
	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/izap"
)

var configFile string
var port string
var mode string
var StartCmd = &cobra.Command{
	Use:          "server",
	Short:        "Start API server",
	Example:      "go-admin server -c config.yaml",
	SilenceUsage: true,
	PreRun:       setup,
	RunE:         run,
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func setup(cmd *cobra.Command, args []string) {
	viper.BindPFlags(cmd.Flags()) // nolint: errcheck
	// viper.SetEnvPrefix("oam")
	// // OAM_CONFIGFILE
	// viper.BindEnv("config") // nolint: errcheck

	//1. 读取配置
	deployed.SetupConfig(configFile)
	//2. 设置日志
	deployed.SetupLogger()
	//3. 初始化数据库链接
	dao.SetupDatabase(deployed.DbConfig.Driver, deployed.DbConfig.Source, deployed.EnabledDB)
	//4. 接口访问控制加载
	deployed.SetupCasbin()

	izap.Sugar.Info(`starting api server`)
}

func run(cmd *cobra.Command, args []string) error {
	if viper.GetString("mode") == infra.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := router.InitRouter()

	srv := &http.Server{
		Addr:    net.JoinHostPort(deployed.AppConfig.Host, deployed.AppConfig.Port),
		Handler: engine,
	}
	go func() {
		jobs.InitJob()
		jobs.Setup()
	}()

	go func() {
		// 服务连接
		if deployed.SslConfig.Enable {
			if err := srv.ListenAndServeTLS(deployed.SslConfig.Pem, deployed.SslConfig.KeyStr); err != nil && err != http.ErrServerClosed {
				izap.Sugar.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				izap.Sugar.Fatal("listen: ", err)
			}
		}
	}()
	tip()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", time.Now().Format("2006-01-02 15:04:05"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	fmt.Println("Server exiting")

	return nil
}

func tip() {
	content, _ := ioutil.ReadFile("./static/go-admin.txt")
	fmt.Println(textcolor.Red(string(content)))
	fmt.Printf("%s \n\n", `欢迎使用 `+textcolor.Green(`go-admin `+builder.Version)+` 可以使用 `+textcolor.Red(`-h`)+` 查看命令`)
	fmt.Println(textcolor.Green("Server run at:"))
	fmt.Printf("-  Local:   http://localhost:%s/ \r\n", deployed.AppConfig.Port)
	fmt.Printf("-  Network: http://%s:%s/ \r\n", infra.LanIP(), deployed.AppConfig.Port)
	fmt.Println(textcolor.Green("Swagger run at:"))
	fmt.Printf("-  Local:   http://localhost:%s/swagger/index.html \r\n", deployed.AppConfig.Port)
	fmt.Printf("-  Network: http://%s:%s/swagger/index.html \r\n", infra.LanIP(), deployed.AppConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", time.Now().Format("2006-01-02 15:04:05"))
}
