package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/thinkgos/sharp/builder"
)

var StartCmd = &cobra.Command{
	Use:     "version",
	Short:   "Get version info",
	Example: "go-admin version",
	RunE:    run,
}

func run(*cobra.Command, []string) error {
	printVersion()
	return nil
}

func printVersion() {
	fmt.Printf("Model: %s\r\n", builder.Model)
	fmt.Printf("Version: %s\r\n", builder.Version)
	fmt.Printf("API version: %s\r\n", builder.APIVersion)
	fmt.Printf("Go version: %s\r\n", runtime.Version())
	fmt.Printf("Git commit: %s\r\n", builder.GitCommit)
	fmt.Printf("Git full commit: %s\r\n", builder.GitFullCommit)
	fmt.Printf("Build time: %s\r\n", builder.BuildTime)
	fmt.Printf("OS/Arch: %s/%s\r\n", runtime.GOOS, runtime.GOARCH)
}
