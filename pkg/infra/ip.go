package infra

import (
	"errors"
	"net"
)

func ActiveIPNet() (*net.IPNet, error) {
	ifcs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, ifc := range ifcs {
		if (ifc.Flags & net.FlagUp) != 0 {
			addrs, _ := ifc.Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet, nil
					}
				}
			}
		}
	}
	return nil, errors.New("not any active net interfaces")
}

// LanIP 局域网ip地址
func LanIP() string {
	an, err := ActiveIPNet()
	if err != nil {
		return ""
	}
	return an.IP.String()
}
