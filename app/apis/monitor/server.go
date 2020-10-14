package monitor

import (
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/x-tardis/go-admin/pkg/infra"
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
type Mem struct {
	Total float64 `json:"total,string"`
	Used  float64 `json:"used,string"`
	Free  float64 `json:"free,string"`
	Usage float64 `json:"usage,string"`
}

type Cpu struct {
	CpuInfo []cpu.InfoStat `json:"cpuInfo"`
	Percent float64        `json:"percent,string"`
	CpuNum  int            `json:"cpuNum"`
}

type Disk struct {
	Total float64 `json:"total,string"`
	Free  float64 `json:"free,string"`
}

type SystemInfo struct {
	Code int  `json:"code"`
	Os   Os   `json:"os"`
	Mem  Mem  `json:"mem"`
	Cpu  Cpu  `json:"cpu"`
	Disk Disk `json:"disk"`
}

// @Summary 系统信息
// @Description 获取JSON
// @Tags 系统信息
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/setting/serverInfo [get]
func ServerInfo(c *gin.Context) {
	projectDir, _ := os.Getwd()
	dis, _ := disk.Usage("/")
	vmem, _ := mem.VirtualMemory()
	percent, _ := cpu.Percent(0, false)
	cpuInfo, _ := cpu.Info()
	cpuNum, _ := cpu.Counts(false)

	c.JSON(http.StatusOK, SystemInfo{
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
			infra.Round(float64(vmem.Total)/GB, 2),
			infra.Round(float64(vmem.Used)/GB, 2),
			infra.Round(float64(vmem.Free)/GB, 2),
			infra.Round(vmem.UsedPercent, 2),
		},
		Cpu{
			cpuInfo,
			infra.Round(percent[0], 2),
			cpuNum,
		},
		Disk{
			infra.Round(float64(dis.Total)/GB, 2),
			infra.Round(float64(dis.Free)/GB, 2),
		},
	})
}
