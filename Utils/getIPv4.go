package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/gopacket/pcap"
)

//GetIPv4 获取分配给该网卡的内网ipv4地址
func GetIPv4(device string) (net.IP, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Println("pcap.FindAllDevs failed,err:", err)
		return nil, err
	}
	for i := range devices {
		if devices[i].Name != device {
			continue
		}
		for _, addr := range devices[i].Addresses {
			ipv4 := addr.IP.To4()
			if ipv4 == nil {
				continue
			}
			return ipv4, nil
		}
	}
	return nil, errors.New("don't have ipv4 address")
}

//GetDefaultOptions 获取默认扫描范围和网关
func GetDefaultOptions() (ifname string, scanRange string, gateway string, err error) {
	var myIP net.IP
	//获取所有设备
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
