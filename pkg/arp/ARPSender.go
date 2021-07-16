package arp

import (
	"ARPSpoofing/models"
	"ARPSpoofing/utils"
	"fmt"
	"net"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//ARPSender ARP广播包发送器
type ARPSender struct {
	ethDstMAC net.HardwareAddr
	ethSrcMAC net.HardwareAddr
	arpDstMAC net.HardwareAddr
	arpSrcMAC net.HardwareAddr
	arpSrcIP  net.IP
	handle    *pcap.Handle
}

func newARPSender(iface *net.Interface) (models.Sender, error) {
	handle, err := pcap.OpenLive(iface.Name, 65535, true, 1000)
	if err != nil {
		return nil, fmt.Errorf("OpenLive Error:%v", err)
	}
	myIP, err := utils.GetIPv4ByIface(iface)
	if err != nil {
		return nil, fmt.Errorf("IPv4 error:%v", err)
	}
	ethDstMAC, err := net.ParseMAC("ff:ff:ff:ff:ff:ff")
	if err != nil {
		return nil, err
	}
	arpDstMAC, err := net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		return nil, err
	}
	//构造广播包
	return &ARPSender{
		ethSrcMAC: iface.HardwareAddr,
		ethDstMAC: ethDstMAC,
		arpSrcMAC: iface.HardwareAddr,
		arpDstMAC: arpDstMAC,
		arpSrcIP:  myIP,
		handle:    handle,
	}, nil
}

//Send 发送数据包
func (s *ARPSender) Send(dstIP net.IP) error {
	return SendARP(s.handle, s.ethDstMAC, s.ethSrcMAC, s.arpDstMAC, s.arpSrcMAC, dstIP, s.arpSrcIP, layers.ARPRequest)
}

//Recv 接收数据包
func (s *ARPSender) Recv(out chan *models.Host) error {
	ReceiveARP(s.handle, s.ethSrcMAC, out, true)
	return nil
}
