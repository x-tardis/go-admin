package system

import (
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/thinkgos/go-core-package/extmath"

	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// Os os信息
type Os struct {
	GoOs         string `json:"goOs"`
	Arch         string `json:"arch"`
	Mem          int    `json:"mem"`
	Compiler     string `json:"compiler"`
	Version      string `json:"version"`
	NumGoroutine int    `json:"numGoroutine"`
	Ip           string `json:"ip"`
	ProjectDir   string `json:"projectDir"`
}

// Mem mem信息
type Mem struct {
	Total       float64 `json:"total,string"`
	Used        float64 `json:"used,string"`
	Free        float64 `json:"free,string"`
	UsedPercent float64 `json:"usedPercent,string"`
}

// Cpu cpu信息
type Cpu struct {
	CpuInfo []cpu.InfoStat `json:"cpuInfo"`
	Percent float64        `json:"percent,string"`
	CpuNum  int            `json:"cpuNum"`
}

// Disk disk信息
type Disk struct {
	Total       float64 `json:"total,string"`
	Used        float64 `json:"used,string"`
	Free        float64 `json:"free,string"`
	UsedPercent float64 `json:"usedPercent,string"`
}

// SystemInfos system info
type SystemInfos struct {
	Code int  `json:"code"`
	Os   Os   `json:"os"`
	Mem  Mem  `json:"mem"`
	Cpu  Cpu  `json:"cpu"`
	Disk Disk `json:"disk"`
}

// @tags 系统信息
// @summary 系统信息
// @description 获取系统信息
// @accept json
// @produce json
// @success 200 {object} SystemInfos "成功回复"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 404 {object} servers.Response "未找到相关信息"
// @failure 417 {object} servers.Response "客户端请求头错误"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/system/info [get]
func SystemInfo(c *gin.Context) {
	projectDir, _ := os.Getwd()
	dis, _ := disk.Usage("/")
	vMem, _ := mem.VirtualMemory()
	percent, _ := cpu.Percent(0, false)
	cpuInfo, _ := cpu.Info()
	cpuNum, _ := cpu.Counts(false)

	c.JSON(http.StatusOK, SystemInfos{
		http.StatusOK,
		Os{
			runtime.GOOS,
			runtime.GOARCH,
			runtime.MemProfileRate,
			runtime.Compiler,
			runtime.Version(),
			runtime.NumGoroutine(),
			infra.LanIP(),
			projectDir,
		},
		Mem{
			extmath.Round(float64(vMem.Total)/GB, 2),
			extmath.Round(float64(vMem.Used)/GB, 2),
			extmath.Round(float64(vMem.Free)/GB, 2),
			extmath.Round(vMem.UsedPercent, 2),
		},
		Cpu{
			cpuInfo,
			extmath.Round(percent[0], 2),
			cpuNum,
		},
		Disk{
			extmath.Round(float64(dis.Total)/GB, 2),
			extmath.Round(float64(dis.Used)/GB, 2),
			extmath.Round(float64(dis.Free)/GB, 2),
			extmath.Round(dis.UsedPercent, 2),
		},
	})
}

const INDEX = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>GO-ADMIN欢迎您</title>
<style>
body{
  margin:0; 
  padding:0; 
  overflow-y:hidden
}
</style>
<script src="http://libs.baidu.com/jquery/1.9.0/jquery.js"></script>
<script type="text/javascript"> 
window.onerror=function(){return true;} 
$(function(){ 
  headerH = 0;  
  var h=$(window).height();
  $("#iframe").height((h-headerH)+"px"); 
});
</script>
</head>
<body>
<iframe id="iframe" frameborder="0" src="http://doc.zhangwj.com/github.com/x-tardis/go-admin-site/" style="width:100%;"></iframe>
</body>
</html>
`

func HelloWorld(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, INDEX)
}

// @tags 系统信息
// @summary ping/pong test
// @description  ping/pong test
// @accept json
// @produce json
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 404 {object} servers.Response "未找到相关信息"
// @failure 417 {object} servers.Response "客户端请求头错误"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/system/ping [get]
func Ping(c *gin.Context) {
	servers.JSON(c, http.StatusOK, servers.WithMsg("pong"))
}
