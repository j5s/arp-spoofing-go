package utils

import (
	"log"

	"github.com/google/gopacket/pcap"
)

//GetDefaultIfname 获取默认网卡名
func GetDefaultIfname() string {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Println("pcap.FindAllDevs failed,err:", err)
		return ""
	}
	if len(devices) == 0 {
		return ""
	}
	return devices[0].Name
}
