package util

import (
	"net"
)

const LocalHost = "127.0.0.1"

func LocalIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return LocalHost
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return LocalHost
}
