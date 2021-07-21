package assembly

import (
	"ARPSpoofing/settings"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//Run 启动http监听程序
func Run() {
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
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetCh := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packetCh:
			if packet.NetworkLayer() == nil {
				continue
			}
			if packet.TransportLayer() == nil {
				continue
			}
			if packet.TransportLayer() != layers.LayerTypeTCP {
				continue
			}

		case <-ticker:

		}
	}
}
