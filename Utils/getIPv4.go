package utils

import (
	"errors"
	"net"
)

//获取分配给该网卡的内网ipv4地址
func GetIPv4ByIface(iface *net.Interface) (net.IP, error) {
	addrs, err := iface.Addrs() //ipv6/mask,ipv4/mask
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ipmask, ok := addr.(*net.IPNet)
		if ok {
			ipv4 := ipmask.IP.To4()
			if ipv4 != nil {
				return ipv4, nil
			}
		}
	}
	return nil, errors.New("don't have ipv4 address")
}
