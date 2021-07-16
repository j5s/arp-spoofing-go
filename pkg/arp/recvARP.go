package arp

import (
	"ARPSpoofing/models"
	"bytes"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//ReceiveARP 接收ARP包
func ReceiveARP(handle *pcap.Handle, myMAC net.HardwareAddr, out chan *models.Host, isRecvAll bool) {
	//1.创建数据源
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	//2.从数据源中读数据
	for packet := range src.Packets() {
		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer == nil { //不是ARP包
			continue
		}
		arp := arpLayer.(*layers.ARP)
		if isRecvAll {
			// log.Println("收到数据包:", net.IP(arp.SourceProtAddress).String(), net.HardwareAddr(arp.SourceHwAddress).String())
			// log.Println("收到数据包:", net.IP(arp.DstProtAddress).String(), net.HardwareAddr(arp.DstHwAddress).String())
			out <- &models.Host{
				IP:  net.IP(arp.SourceProtAddress).String(),
				MAC: net.HardwareAddr(arp.SourceHwAddress).String(),
			}
			out <- &models.Host{
				IP:  net.IP(arp.DstProtAddress).String(),
				MAC: net.HardwareAddr(arp.DstHwAddress).String(),
			}
		} else {
			//不是ARP响应包
			if arp.Operation != layers.ARPReply {
				continue
			}
			//不是发送给我的包
			if false == bytes.Equal([]byte(myMAC), arp.DstHwAddress) {
				continue
			}
			out <- &models.Host{
				IP:  net.IP(arp.SourceProtAddress).String(),
				MAC: net.HardwareAddr(arp.SourceHwAddress).String(),
			}
		}
	}
	log.Println("ARPRecive quit")
}
