package misc

import (
	"fmt"
	"log"
	"runtime"

	"github.com/thinkgos/sharp/builder"
)

func PrintVersion() {
	fmt.Printf("Model: %s\r\n", builder.Model)
	fmt.Printf("Version: %s\r\n", builder.Version)
	fmt.Printf("API version: %s\r\n", builder.APIVersion)
	fmt.Printf("Go version: %s\r\n", runtime.Version())
	fmt.Printf("Git commit: %s\r\n", builder.GitCommit)
	fmt.Printf("Git full commit: %s\r\n", builder.GitFullCommit)
	fmt.Printf("Build time: %s\r\n", builder.BuildTime)
	fmt.Printf("OS/Arch: %s/%s\r\n", runtime.GOOS, runtime.GOARCH)
}

func HandlerError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
