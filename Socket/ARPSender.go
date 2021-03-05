package Socket

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"ARPSpoofing/Utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type ARPSender struct {
	ethDstMAC net.HardwareAddr
	ethSrcMAC net.HardwareAddr
	arpDstMAC net.HardwareAddr
	arpSrcMAC net.HardwareAddr
	arpSrcIP  net.IP
	handle    *pcap.Handle
}

func newARPSender(iface *net.Interface) (Sender, error) {
	log.Println("new ARP Sender")
	handle, err := pcap.OpenLive(iface.Name, 65535, true, 1000)
	if err != nil {
		return nil, fmt.Errorf("OpenLive Error:%v", err)
	}
	myIP, err := Utils.GetIPv4ByIface(iface)
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
	return &ARPSender{
		ethSrcMAC: iface.HardwareAddr,
		ethDstMAC: ethDstMAC,
		arpSrcMAC: iface.HardwareAddr,
		arpDstMAC: arpDstMAC,
		arpSrcIP:  myIP,
		handle:    handle,
	}, nil
}
func (this *ARPSender) Send(dstIP net.IP) error {
	return SendARP(this.handle, this.ethDstMAC, this.ethSrcMAC, this.arpDstMAC, this.arpSrcMAC, dstIP, this.arpSrcIP, layers.ARPRequest)
}
func (this *ARPSender) Recv(out chan *HostItem) error {
	ReceiveARP(this.handle, this.ethSrcMAC, out, true)
	return nil
}

func SendARP(handle *pcap.Handle, ethDstMAC, ethSrcMAC, arpDstMAC, arpSrcMAC net.HardwareAddr, arpDstIP, arpSrcIP net.IP, opt uint16) error {
	//构造ARP数据包
	eth := layers.Ethernet{
		SrcMAC:       ethSrcMAC,
		DstMAC:       ethDstMAC,
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         opt, //request or reply
		SourceHwAddress:   arpSrcMAC,
		SourceProtAddress: arpSrcIP.To4(),
		DstHwAddress:      arpDstMAC,
		DstProtAddress:    arpDstIP.To4(),
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
		return fmt.Errorf("序列化失败:%s", err)
	}
	//发送ARP数据包
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return fmt.Errorf("WritePacketData失败:%s", err)
	}
	return nil
}

func ReceiveARP(handle *pcap.Handle, myMAC net.HardwareAddr, out chan *HostItem, isRecvAll bool) {
	log.Println("RecvARP start")
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	for packet := range src.Packets() {
		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer == nil { //不是ARP包
			continue
		}
		arp := arpLayer.(*layers.ARP)
		if isRecvAll {
			// log.Println("收到数据包:", net.IP(arp.SourceProtAddress).String(), net.HardwareAddr(arp.SourceHwAddress).String())
			// log.Println("收到数据包:", net.IP(arp.DstProtAddress).String(), net.HardwareAddr(arp.DstHwAddress).String())
			out <- NewHostItem(net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))
			out <- NewHostItem(net.IP(arp.DstProtAddress), net.HardwareAddr(arp.DstHwAddress))
		} else {
			if arp.Operation != layers.ARPReply {
				continue
			}
			if !bytes.Equal([]byte(myMAC), arp.DstHwAddress) {
				continue
			}
			out <- NewHostItem(net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))
		}
	}
	log.Println("ARPRecive quit")
}
