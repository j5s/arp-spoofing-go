package arp

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

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
