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
	"github.com/google/gops/agent"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thinkgos/go-core-package/lib/ternary"
	"github.com/thinkgos/go-core-package/lib/textcolor"
	"github.com/thinkgos/sharp/builder"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/misc"
	"github.com/x-tardis/go-admin/pkg/bcron"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/jobs"
	"github.com/x-tardis/go-admin/pkg/tip"
	"github.com/x-tardis/go-admin/router"
)

var configFile string
var port string
var mode string
var Cmd = &cobra.Command{
	Use:          "server",
	Short:        "Start API server",
	Example:      fmt.Sprintf("%s server -c config/config.yaml", builder.Name),
	SilenceUsage: true,
	PreRun:       setup,
	RunE:         run,
	PostRun:      postRun,
}

func init() {
	Cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yaml", "Start server with provided configuration file")
	Cmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "Tcp port server listening on")
	Cmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,debug,prod")
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
	dao.SetupDatabase(dao.DbConfig)
	// 4. 接口访问控制加载
	deployed.SetupCasbin()
	// 时钟定时器
	bcron.Cron.Start()
}

func run(cmd *cobra.Command, args []string) error {
	var err error

	fmt.Println(textcolor.Red("starting server..."))

	go func() {
		misc.HandlerError(agent.Listen(deployed.ViperGops()))
		time.Sleep(time.Millisecond * 100)
		jobs.Startup()
	}()
	// 设置gin的工作模式
	gin.SetMode(ternary.IfString(deployed.IsModeProd(), gin.ReleaseMode, gin.DebugMode))

	engine := router.InitRouter()
	addr := net.JoinHostPort(deployed.AppConfig.Host, deployed.AppConfig.Port)

	showTip()

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
		log.Fatal("listen and serve: ", err)
	}
	return nil
}

func postRun(*cobra.Command, []string) {
	bcron.Cron.Stop()
}

func showTip() {
	t := tip.Tip{
		textcolor.Red(infra.Banner),
		textcolor.Green(deployed.AppConfig.Name),
		textcolor.Magenta(builder.Version),
		textcolor.Magenta("-h"),
		textcolor.Green("Server run at:"),
		textcolor.Green("Swagger run at:"),
		infra.LanIP(),
		deployed.AppConfig.Port,
		textcolor.Green("Server run on PID:"),
		textcolor.Red(cast.ToString(os.Getpid())),
		textcolor.Magenta("Control + C"),
	}
	tip.Show(t)
}
