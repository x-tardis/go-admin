package monitor

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/tools"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// @Summary 系统信息
// @Description 获取JSON
// @Tags 系统信息
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/settings/serverInfo [get]
func ServerInfo(c *gin.Context) {
	osDic := make(map[string]interface{})
	osDic["goOs"] = runtime.GOOS
	osDic["arch"] = runtime.GOARCH
	osDic["mem"] = runtime.MemProfileRate
	osDic["compiler"] = runtime.Compiler
	osDic["version"] = runtime.Version()
	osDic["numGoroutine"] = runtime.NumGoroutine()
	osDic["ip"] = infra.LanIP()
	osDic["projectDir"] = tools.GetCurrentPath()

	diskDic := make(map[string]interface{})
	dis, _ := disk.Usage("/")
	diskDic["total"] = int(dis.Total) / GB
	diskDic["free"] = int(dis.Free) / GB

	memDic := make(map[string]interface{})
	_mem, _ := mem.VirtualMemory()
	memDic["total"] = int(_mem.Total) / GB
	memDic["used"] = int(_mem.Used) / GB
	memDic["free"] = int(_mem.Free) / GB
	memDic["usage"] = int(_mem.UsedPercent)

	cpuDic := make(map[string]interface{})
	cpuDic["cpuInfo"], _ = cpu.Info()
	percent, _ := cpu.Percent(0, false)
	cpuDic["Percent"] = infra.Round(percent[0], 2)
	cpuDic["cpuNum"], _ = cpu.Counts(false)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"os":   osDic,
		"mem":  memDic,
		"cpu":  cpuDic,
		"disk": diskDic,
	})
}
