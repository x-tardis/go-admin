package version

import (
	"github.com/spf13/cobra"

	"github.com/x-tardis/go-admin/misc"
)

var StartCmd = &cobra.Command{
	Use:     "version",
	Short:   "Get version info",
	Example: "go-admin version",
	RunE: func(*cobra.Command, []string) error {
		misc.PrintVersion()
		return nil
	},
}
