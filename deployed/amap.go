package deployed

import (
	"github.com/thinkgos/go-core-package/extnet"

	"github.com/x-tardis/go-admin/pkg/amap"
)

var defaultAmap = amap.New("3fabc36c20379fbb9300c79b19d5d05e")

// 获取外网ip地址
func IPLocation(ip string) string {
	if extnet.IsIntranet(ip) {
		return "intranet location"
	}
	rsp, err := defaultAmap.IPLocation(ip)
	if err != nil || (rsp.Province == "" && rsp.City == "") {
		return "unknown location"
	}

	if rsp.Province == "" {
		return rsp.City
	}
	return rsp.Province + "-" + rsp.City
}
