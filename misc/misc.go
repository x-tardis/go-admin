package misc

import (
	"log"
	"os"
	"runtime"
	"text/template"

	"github.com/thinkgos/sharp/builder"
)

const version = `  Model:            {{.Model}}
  Version:          {{.Version}}
  API version:      {{.APIVersion}}
  Go version:       {{.GoVersion}}
  Git commit:       {{.GitCommit}}
  Git full commit:  {{.GitFullCommit}}
  Build time:       {{.BuildTime}}
  OS/Arch:          {{.GOOS}}/{{.GOARCH}}
`

type Version struct {
	Model         string
	Version       string
	APIVersion    string
	GoVersion     string
	GitCommit     string
	GitFullCommit string
	BuildTime     string
	GOOS          string
	GOARCH        string
}

func PrintVersion() {
	v := Version{
		builder.Model,
		builder.Version,
		builder.APIVersion,
		runtime.Version(),
		builder.GitCommit,
		builder.GitFullCommit,
		builder.BuildTime,
		runtime.GOOS, runtime.GOARCH,
	}
	template.Must(template.New("version").Parse(version)).Execute(os.Stdout, v)
}

func HandlerError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
