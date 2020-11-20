package tip

import (
	"html/template"
	"os"
	"runtime"

	"github.com/thinkgos/sharp/builder"
)

const versionTpl = `  Name:             {{.Name}}
  Model:            {{.Model}}
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
	Name          string
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
		builder.Name,
		builder.Model,
		builder.Version,
		builder.APIVersion,
		runtime.Version(),
		builder.GitCommit,
		builder.GitFullCommit,
		builder.BuildTime,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
	}
	template.Must(template.New("version").Parse(versionTpl)).
		Execute(os.Stdout, v) // nolint: errcheck
}
