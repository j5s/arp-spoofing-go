package arp

import (
	"ARPSpoofing/debug"
	"ARPSpoofing/models"
	"ARPSpoofing/pkg/routine"
	"bytes"
	"context"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//ReceiveARP 接收ARP包
func ReceiveARP(ctx context.Context, handle *pcap.Handle, myMAC net.HardwareAddr, isRecvAll bool) chan *models.Host {
	out := make(chan *models.Host, 1024)
	go func(ctx context.Context) {
		debug.Println("启动了一个接收ARP报文的协程:", routine.GetGID())
		defer func() {
			close(out)
			debug.Println("接收ARP报文的线程退出:", routine.GetGID())
		}()
		//1.创建数据源
		src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
		//2.从数据源中读数据
		for packet := range src.Packets() {
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil { //不是ARP包
				continue
			}
			arp, ok := arpLayer.(*layers.ARP)
			if !ok {
				continue
			}
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
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}(ctx)
	return out
}
