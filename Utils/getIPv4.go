package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

//GetIPv4ByIface 获取分配给该网卡的内网ipv4地址
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

//GetDefaultScanRangeGateway 获取默认扫描范围和网关
func GetDefaultScanRangeGateway(ifname string) (scanRange string, gateway string, err error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return "", "", err
	}
	myIP, err := GetIPv4ByIface(iface)
	if err != nil {
		return "", "", err
	}
	scanRange = fmt.Sprintf("%s/24", myIP.String())
	gateway = fmt.Sprintf("%s.1", strings.Join(strings.Split(myIP.String(), ".")[:3], "."))
	return scanRange, gateway, nil
}
