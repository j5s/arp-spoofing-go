package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const UDPHeaderLen int = 8

func SendUDP(handle *pcap.Handle, dstMAC, srcMAC net.HardwareAddr, dstIPstr, srcIPstr string, dstPort, srcPort layers.UDPPort, dataLen int) error {
	dstIP := net.ParseIP(dstIPstr)
	if dstIP == nil {
		return errors.New("错误的IP地址")
	}
	srcIP := net.ParseIP(srcIPstr)
	if srcIP == nil {
		return errors.New("错误的IP地址")
	}

	Eth := layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}
	IP := layers.IPv4{
		Version:  4,
		TTL:      255,
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolUDP,
	}
	UDP := layers.UDP{
		SrcPort: srcPort,
		DstPort: dstPort,
		Length:  uint16(UDPHeaderLen + dataLen),
	}
	UDP.SetNetworkLayerForChecksum(&IP) //计算UDP效验和，因为要构造伪首部，所以需要传入IP层的数据
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, &Eth, &IP, &UDP); err != nil {
		return fmt.Errorf("序列化发生错误:%v", err)
	}
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return fmt.Errorf("发送数据包出现错误:%v", err)
	}
	return nil
}

func RecvICMP(handle *pcap.Handle, out chan *HostItem) {
	log.Println("RecvICMP start")
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	for packet := range src.Packets() {
		_ICMPLayerV4 := packet.Layer(layers.LayerTypeICMPv4)
		if _ICMPLayerV4 == nil {
			continue
		}
		ICMPLayerV4 := _ICMPLayerV4.(*layers.ICMPv4)

		fmt.Println(ICMPLayerV4.TypeCode.GoString())

		_EtherLayer := packet.Layer(layers.LayerTypeEthernet)
		if _EtherLayer == nil {
			continue
		}
		EtherLayer := _EtherLayer.(*layers.Ethernet)
		srcMAC := EtherLayer.SrcMAC
		_IPv4Layer := packet.Layer(layers.LayerTypeIPv4)
		if _IPv4Layer == nil {
			continue
		}
		IPv4Layer := _IPv4Layer.(*layers.IPv4)
		srcIP := IPv4Layer.SrcIP
		out <- newHostItem(srcIP, srcMAC)
	}
	log.Println("RecvICMP quit")
}

func SendUDPv2(srcIPstr, dstIPstr string, srcPort, dstPort, dataLen int) error {
	dstIP := net.ParseIP(dstIPstr)
	if dstIP == nil {
		return errors.New("无效的IP地址")
	}
	srcIP := net.ParseIP(srcIPstr)
	if srcIP == nil {
		return errors.New("无效的IP地址")
	}
	conn, err := net.ListenPacket("ip4:udp", srcIP.String()) //创建一个UDP套接字
	if err != nil {
		return fmt.Errorf("ListenPacket失败:%v", err)
	}
	defer conn.Close()
	ip := &layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolUDP,
	}
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(srcPort),
		DstPort: layers.UDPPort(dstPort),
		Length:  uint16(dataLen),
	}
	udp.SetNetworkLayerForChecksum(ip)
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, udp); err != nil {
		return fmt.Errorf("序列化失败:%v", err)
	}
	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstIP}); err != nil {
		return fmt.Errorf("发送数据包失败:%v", err)
	}
	return nil
}
