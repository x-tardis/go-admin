package misc

import (
	"log"
	"os"
	"runtime"
	"text/template"

	"github.com/thinkgos/sharp/builder"
)

const versionTpl = `  Model:            {{.Model}}
  Version:          {{.Version}}
  API version:      {{.APIVersion}}
  Go version:       {{.GoVersion}}
  Git commit:       {{.GitCommit}}
  Git full commit:  {{.GitFullCommit}}
  Build time:       {{.BuildTime}}
  OS/Arch:          {{.GOOS}}/{{.GOARCH}}
  NumCPU:           {{.NumCPU}}
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
	NumCPU        int
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
		runtime.NumCPU(),
	}
	template.Must(template.New("version").Parse(versionTpl)).Execute(os.Stdout, v)
}

func HandlerError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
