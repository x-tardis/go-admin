package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thinkgos/x/builder"

	"github.com/x-tardis/go-admin/app/cmd/config"
	"github.com/x-tardis/go-admin/app/cmd/migrate"
	"github.com/x-tardis/go-admin/app/cmd/server"
	"github.com/x-tardis/go-admin/app/cmd/version"
)

func init() {
	rootCmd.AddCommand(
		server.Cmd,
		migrate.Cmd,
		version.Cmd,
		config.Cmd,
	)
}

var rootCmd = &cobra.Command{
	Use:          builder.Name,
	Short:        builder.Name,
	SilenceUsage: true,
	Long:         builder.Name,
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
