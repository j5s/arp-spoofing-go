package arp

import (
	"fmt"
	"log"
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
		log.Println("gopacket.SerializeLayers failed,err:", err)
		return err
	}
	//发送ARP数据包
	err := handle.WritePacketData(buf.Bytes())
	if err != nil {
		fmt.Println("handle.WritePacketData failed,err:", err)
		return err
	}
	return nil
}
