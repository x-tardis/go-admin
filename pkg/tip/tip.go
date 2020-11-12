package tip

import (
	"os"
	"text/template"
)

const tipText = `  {{.Banner}}

欢迎使用 {{.Name}} {{.Version}} 可以使用 {{.H}} 查看命令
{{.ServerTitle}}:
	-  Local:   http://localhost:{{.Port}}
	-  Network: http://{{.ExtranetIP}}:{{.Port}}
{{.SwaggerTitle}}:
	-  Local:   http://localhost:{{.Port}}/swagger/index.html
	-  Network: http://{{.ExtranetIP}}:{{.Port}}/swagger/index.html
{{.PidTitle}}: {{.PID}}
Enter {{.Kill}} Shutdown Server

`

// Tip tip
type Tip struct {
	Banner       string // 横幅
	Name         string // 应用名称
	Version      string // 应用版本
	H            string // 一般为 -h
	ServerTitle  string // 服务标题
	SwaggerTitle string // swagger标题
	ExtranetIP   string // 外网ip地址
	Port         string // 端口
	PidTitle     string // pid标题
	PID          string // pid
	Kill         string // 一般为 Control + C
}

// Show show tip
func Show(t Tip) {
	template.Must(template.New("tip").Parse(tipText)).
		Execute(os.Stdout, t) // nolint: errcheck
}
