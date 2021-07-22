package assembly

import (
	"ARPSpoofing/settings"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

//Run 启动http监听程序
func Run(ctx context.Context) {
	ifaceName, err := settings.Options.Get("ifname")
	if err != nil {
		log.Println("请先选择网卡")
		return
	}
	handle, err := pcap.OpenLive(ifaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Println("pcap.OpenLive failed,err:", err)
		return
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetCh := packetSource.Packets()
	ticker := time.Tick(time.Minute)

	streamFactory := &httpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)
	for {
		select {
		case packet := <-packetCh:
			if packet == nil {
				return
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP || packet.ApplicationLayer() == nil {
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			//将TCP包重组到流中
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-ticker:
			assembler.FlushOlderThan(time.Now().Add(time.Second * -20))
		case <-ctx.Done():
			fmt.Printf("\r[*] 嗅探协程退出成功\n")
			return
		}
	}
}
