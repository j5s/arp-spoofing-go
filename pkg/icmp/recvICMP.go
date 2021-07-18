package icmp

import (
	"ARPSpoofing/debug"
	"ARPSpoofing/models"
	"ARPSpoofing/pkg/routine"
	"context"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//RecvICMP 接收ICMP数据包
func RecvICMP(ctx context.Context, handle *pcap.Handle) <-chan *models.Host {
	outCh := make(chan *models.Host, 2048)
	go func(ctx context.Context) {
		debug.Println("启动了一个接收ICMP报文的协程:", routine.GetGID())
		defer func() {
			close(outCh)
			debug.Println("接收ICMP的协程退出：", routine.GetGID())
		}()
		src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
		for packet := range src.Packets() {
			//1.ICMP层
			_ICMPLayerV4 := packet.Layer(layers.LayerTypeICMPv4)
			if _ICMPLayerV4 == nil {
				continue
			}
			//2.以太网层
			_EtherLayer := packet.Layer(layers.LayerTypeEthernet)
			if _EtherLayer == nil {
				continue
			}
			EtherLayer, ok := _EtherLayer.(*layers.Ethernet)
			if !ok {
				continue
			}
			//3.IP层
			_IPv4Layer := packet.Layer(layers.LayerTypeIPv4)
			if _IPv4Layer == nil {
				continue
			}
			IPv4Layer, ok := _IPv4Layer.(*layers.IPv4)
			if !ok {
				continue
			}
			//4.送入管道
			outCh <- &models.Host{
				IP:  IPv4Layer.SrcIP.String(),
				MAC: EtherLayer.SrcMAC.String(),
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}(ctx)
	return outCh
}
