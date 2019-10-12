package utils

import (
	"net"
)

// GetIP func is get ipv4 address
func GetIP() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ipv4s []string
	for _, a := range addrs {
		ipaddr, ok := a.(*net.IPNet)
		if ok && !ipaddr.IP.IsLoopback() {
			if ipaddr.IP.To4() != nil {
				ipv4s = append(ipv4s, ipaddr.IP.String())
			}
		}
	}

	return ipv4s, nil
}
