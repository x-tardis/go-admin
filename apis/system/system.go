package system

import (
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/thinkgos/meter"
	"github.com/thinkgos/x/builder"
	"github.com/thinkgos/x/extmath"

	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// Os os信息
type Os struct {
	Os             string `json:"os"`
	Arch           string `json:"arch"`
	NumCPU         int    `json:"numCpu,string"`
	MemProfileRate int    `json:"memProfileRate,string"`
	GoVersion      string `json:"goVersion"`
	Compiler       string `json:"compiler"`
	GoMaxProcs     int    `json:"gomaxprocs,string"`
	Threads        int    `json:"threads,string"`
	NumGoroutine   int    `json:"numGoroutine,string"`
}

// Mem mem信息
type Mem struct {
	Total       string  `json:"total"`
	Used        string  `json:"used"`
	Free        string  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

// Cpu cpu信息
type Cpu struct {
	NumCPU  int            `json:"numCpu,string"`
	Percent float64        `json:"percent"`
	CpuInfo []cpu.InfoStat `json:"cpuInfo"`
}

// Disk disk信息
type Disk struct {
	Total       string  `json:"total"`
	Used        string  `json:"used"`
	Free        string  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

type App struct {
	Model         string `json:"model"`         // 型号
	Ppid          int    `json:"ppid,string"`   // ppid
	Pid           int    `json:"pid,string"`    // pid
	Version       string `json:"version"`       // 版本
	APIVersion    string `json:"apiVersion"`    // api版本
	BuildTime     string `json:"buildTime"`     // 编译时间
	GitCommit     string `json:"gitCommit"`     // git提交码
	GitFullCommit string `json:"gitFullCommit"` // git提交全码
	ProjectDir    string `json:"projectDir"`    // 工作目录
	Ip            string `json:"ip"`
	SetupTime     string `json:"setupTime"`   // 程序启动日期
	RunningTime   string `json:"runningTime"` // 程序运行时间
	Uptime        string `json:"uptime"`      // 系统运行时间
	Mem           Mem    `json:"mem"`
}

// SystemInfos system info
type SystemInfos struct {
	Code int  `json:"code"`
	Os   Os   `json:"os"`
	Mem  Mem  `json:"mem"`
	Cpu  Cpu  `json:"cpu"`
	Disk Disk `json:"disk"`
	App  App  `json:"app"`
}

var setupTime = time.Now()

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
	cpuNum, _ := cpu.Counts(true)

	sysInfo := syscall.Sysinfo_t{} // Uptime = 秒
	syscall.Sysinfo(&sysInfo)      // nolint: errcheck

	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	c.JSON(http.StatusOK, SystemInfos{
		http.StatusOK,
		Os{
			runtime.GOOS,
			runtime.GOARCH,
			runtime.NumCPU(),
			runtime.MemProfileRate,
			runtime.Version(),
			runtime.Compiler,
			runtime.GOMAXPROCS(0),
			pprof.Lookup("threadcreate").Count(),
			runtime.NumGoroutine(),
		},
		Mem{
			meter.ByteSize(vMem.Total).String(),
			meter.ByteSize(vMem.Used).String(),
			meter.ByteSize(vMem.Free).String(),
			extmath.Round(vMem.UsedPercent, 2),
		},
		Cpu{
			cpuNum,
			extmath.Round(percent[0], 2),
			cpuInfo,
		},
		Disk{
			meter.ByteSize(dis.Total).String(),
			meter.ByteSize(dis.Used).String(),
			meter.ByteSize(dis.Free).String(),
			extmath.Round(dis.UsedPercent, 2),
		},
		App{
			builder.Model,
			os.Getppid(),
			os.Getpid(),
			builder.APIVersion,
			builder.Version,
			builder.BuildTime,
			builder.GitCommit,
			builder.GitFullCommit,
			projectDir,
			infra.LanIP(),
			setupTime.Format("2006-01-02 15:04:05 Z07:00"),
			time.Since(setupTime).Round(time.Second).String(),
			(time.Duration(sysInfo.Uptime) * time.Second).String(),
			Mem{
				meter.ByteSize(memStats.HeapSys).String(),
				meter.ByteSize(memStats.HeapAlloc).String(),
				meter.ByteSize(memStats.HeapSys - memStats.HeapAlloc).String(),
				extmath.Round(float64(memStats.HeapAlloc)*100/float64(memStats.HeapSys), 2),
			},
		},
	})
}

const Index = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>欢迎您</title>
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

</body>
</html>
`

// <iframe id="iframe" frameborder="0" src="https://doc.zhangwj.com/" style="width:100%;"></iframe>
func HelloWorld(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, Index)
}

// @tags 系统信息
// @summary ping/pong test
// @description  ping/pong test
// @accept json
// @produce json
// @success 200 {object} servers.Response "成功"
// @router /api/v1/system/ping [get]
func Ping(c *gin.Context) {
	servers.OK(c, servers.WithMsg("pong"))
}
