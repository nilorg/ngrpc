package ngrpc

import (
	"errors"
	"net"
)

// LocalIPv4 本地IP
func LocalIPv4(name ...string) (ipAddr string, err error) {
	var addrs []net.Addr
	if len(name) > 0 && name[0] != "" {
		var ief *net.Interface
		ief, err = net.InterfaceByName(name[0])
		if err != nil {
			return
		}
		addrs, err = ief.Addrs()
	} else {
		addrs, err = net.InterfaceAddrs()
	}
	if err != nil {
		return
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddr = ipnet.IP.String()
				return
			}
		}
	}
	err = errors.New("not found local ip")
	return
}
