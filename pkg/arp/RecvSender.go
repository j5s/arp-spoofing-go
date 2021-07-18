package arp

import (
	"ARPSpoofing/models"
	"ARPSpoofing/utils"
	"context"
	"fmt"
	"net"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//RecvSender ARP广播包收发器
type RecvSender struct {
	ethDstMAC net.HardwareAddr
	ethSrcMAC net.HardwareAddr
	arpDstMAC net.HardwareAddr
	arpSrcMAC net.HardwareAddr
	arpSrcIP  net.IP
	handle    *pcap.Handle
}

//NewRecvSender 构造收发函数
func NewRecvSender(iface *net.Interface) (models.RecvSender, error) {
	handle, err := pcap.OpenLive(iface.Name, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("OpenLive Error:%v", err)
	}
	// defer handle.Close() //不能关，还要用
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
	return &RecvSender{
		ethSrcMAC: iface.HardwareAddr,
		ethDstMAC: ethDstMAC,
		arpSrcMAC: iface.HardwareAddr,
		arpDstMAC: arpDstMAC,
		arpSrcIP:  myIP,
		handle:    handle,
	}, nil
}

//Send 发送数据包
func (s *RecvSender) Send(dstIP net.IP) error {
	return SendARP(s.handle, s.ethDstMAC, s.ethSrcMAC, s.arpDstMAC, s.arpSrcMAC, dstIP, s.arpSrcIP, layers.ARPRequest)
}

//Recv 接收数据包
func (s *RecvSender) Recv(ctx context.Context) <-chan *models.Host {
	return ReceiveARP(ctx, s.handle, s.ethSrcMAC, false)
}
