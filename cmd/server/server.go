package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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
	Example:      "go-admin server -c config/config.yaml",
	SilenceUsage: true,
	PreRun:       setup,
	RunE:         run,
	PostRun:      postRun,
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yaml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,debug,prod")
}

func setup(cmd *cobra.Command, args []string) {
	viper.BindPFlags(cmd.Flags()) // nolint: errcheck
	// viper.SetEnvPrefix("oam")
	// // OAM_CONFIGFILE
	// viper.BindEnv("config") // nolint: errcheck

	// 1. 读取配置
	deployed.SetupConfig(configFile)
	// 2. 设置日志
	deployed.SetupLogger()
	// 3. 初始化数据库链接
	dao.SetupDatabase(deployed.DbConfig)
	// 4. 接口访问控制加载
	deployed.SetupCasbin()
}

func run(cmd *cobra.Command, args []string) error {
	var err error

	izap.Sugar.Info(`starting api server`)

	go func() {
		time.Sleep(time.Millisecond * 100)
		jobs.Startup()
	}()

	if viper.GetString("mode") == infra.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := router.InitRouter()
	addr := net.JoinHostPort(deployed.AppConfig.Host, deployed.AppConfig.Port)

	tip()
	// 默认endless服务器会监听下列信号：
	// syscall.SIGHUP，syscall.SIGUSR1，syscall.SIGUSR2，syscall.SIGINT，syscall.SIGTERM和syscall.SIGTSTP
	// 接收到 SIGHUP 信号将触发`fork/restart` 实现优雅重启（kill -1 pid会发送SIGHUP信号）
	// 接收到 syscall.SIGINT或syscall.SIGTERM 信号将触发优雅关机
	// 接收到 SIGUSR2 信号将触发HammerTime
	// SIGUSR1 和 SIGTSTP 被用来触发一些用户自定义的hook函数
	srv := endless.NewServer(addr, engine)
	if deployed.SslConfig.Enable {
		err = srv.ListenAndServeTLS(deployed.SslConfig.Pem, deployed.SslConfig.KeyStr)
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		log.Fatal("listen and serve : ", err)
	}

	return nil
}

func postRun(cmd *cobra.Command, args []string) {
	jobs.Stop()
}

func tip() {
	fmt.Println(textcolor.Red(infra.Banner))
	fmt.Printf("欢迎使用 %s %s 可以使用 %s 查看命令 \n\n", textcolor.Green("go-admin"), textcolor.Magenta(builder.Version), textcolor.Magenta("-h"))
	fmt.Println(textcolor.Green("Server run at:"))
	fmt.Printf("\t-  Local:   http://localhost:%s/ \r\n", deployed.AppConfig.Port)
	fmt.Printf("\t-  Network: http://%s:%s/ \r\n", infra.LanIP(), deployed.AppConfig.Port)
	fmt.Println(textcolor.Green("Swagger run at:"))
	fmt.Printf("\t-  Local:   http://localhost:%s/swagger/index.html \r\n", deployed.AppConfig.Port)
	fmt.Printf("\t-  Network: http://%s:%s/swagger/index.html \r\n", infra.LanIP(), deployed.AppConfig.Port)
	fmt.Printf("%s %s \r\n", textcolor.Green("Server run on PID:"), textcolor.Red(cast.ToString(os.Getpid())))
	log.Printf("Enter %s Shutdown Server\r\n", textcolor.Magenta("Control + C"))
}
