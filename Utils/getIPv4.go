package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/gopacket/pcap"
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

//GetDefaultOptions 获取默认扫描范围和网关
func GetDefaultOptions() (ifname string, scanRange string, gateway string, err error) {
	var myIP net.IP
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Println("pcap.FindAllDevs failed,err:", err)
		return ifname, scanRange, gateway, err
	}
	if len(devices) == 0 {
		return ifname, scanRange, gateway, err
	}
	ifname = devices[0].Name
	for _, addr := range devices[0].Addresses {
		ipv4 := addr.IP.To4()
		if ipv4 != nil {
			myIP = ipv4
			break
		}
	}
	if myIP != nil {
		scanRange = fmt.Sprintf("%s/24", myIP.String())
		gateway = fmt.Sprintf("%s.1", strings.Join(strings.Split(myIP.String(), ".")[:3], "."))
	}
	return ifname, scanRange, gateway, nil
}
