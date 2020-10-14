package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thinkgos/sharp/builder"

	"github.com/x-tardis/go-admin/cmd/config"
	"github.com/x-tardis/go-admin/cmd/migrate"
	"github.com/x-tardis/go-admin/cmd/server"
	"github.com/x-tardis/go-admin/cmd/version"
	"github.com/x-tardis/go-admin/pkg/textcolor"
)

var rootCmd = &cobra.Command{
	Use:          "github.com/x-tardis/go-admin",
	Short:        "github.com/x-tardis/go-admin",
	SilenceUsage: true,
	Long:         `github.com/x-tardis/go-admin`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(textcolor.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎使用 ` + textcolor.Green(`github.com/x-tardis/go-admin `+builder.Version) + ` 可以使用 ` + textcolor.Red(`-h`) + ` 查看命令`
	usageStr1 := `也可以参考 http://doc.zhangwj.com/github.com/x-tardis/go-admin-site/guide/ksks.html 里边的【启动】章节`
	fmt.Printf("%s\n", usageStr)
	fmt.Printf("%s\n", usageStr1)
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
	rootCmd.AddCommand(version.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
