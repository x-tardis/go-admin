package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thinkgos/go-core-package/lib/textcolor"
	"github.com/thinkgos/sharp/builder"

	"github.com/x-tardis/go-admin/cmd/config"
	"github.com/x-tardis/go-admin/cmd/migrate"
	"github.com/x-tardis/go-admin/cmd/server"
	"github.com/x-tardis/go-admin/cmd/version"
)

func init() {
	rootCmd.AddCommand(
		server.StartCmd,
		migrate.StartCmd,
		version.StartCmd,
		config.StartCmd,
	)
}

var rootCmd = &cobra.Command{
	Use:          builder.Name,
	Short:        builder.Name,
	SilenceUsage: true,
	Long:         builder.Name,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip(cmd, args)
			return errors.New(textcolor.Red("requires at least one arg"))
		}
		return nil
	},
	Run: tip,
}

func tip(*cobra.Command, []string) {
	fmt.Printf("欢迎使用 %s %s 可以使用 %s 查看命令\r\n",
		textcolor.Green(builder.Name),
		textcolor.Green(builder.Version),
		textcolor.Red(`-h`))
}

// Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

// AddCommand add command
func AddCommand(cmds ...*cobra.Command) {
	rootCmd.AddCommand(cmds...)
}
