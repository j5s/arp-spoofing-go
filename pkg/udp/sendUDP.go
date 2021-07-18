package udp

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//UDPHeaderLen UDP 头部长度
const UDPHeaderLen int = 8

//SendUDP 发送UDP数据包
func SendUDP(handle *pcap.Handle, dstMAC, srcMAC net.HardwareAddr, dstIPstr, srcIPstr string, dstPort, srcPort layers.UDPPort, dataLen int) error {
	dstIP := net.ParseIP(dstIPstr)
	if dstIP == nil {
		return errors.New("错误的IP地址")
	}
	srcIP := net.ParseIP(srcIPstr)
	if srcIP == nil {
		return errors.New("错误的IP地址")
	}
	//1.构造数据包
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
	//计算UDP效验和，因为要构造伪首部，所以需要传入IP层的数据
	UDP.SetNetworkLayerForChecksum(&IP)
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, &Eth, &IP, &UDP); err != nil {
		return fmt.Errorf("序列化发生错误:%v", err)
	}
	//2.发送数据包
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return fmt.Errorf("发送数据包出现错误:%v", err)
	}
	return nil
}

//SendUDPv2 发送UDP数据包
func SendUDPv2(srcIPstr, dstIPstr string, srcPort, dstPort, dataLen int) error {
	dstIP := net.ParseIP(dstIPstr)
	if dstIP == nil {
		return errors.New("无效的IP地址")
	}
	srcIP := net.ParseIP(srcIPstr)
	if srcIP == nil {
		return errors.New("无效的IP地址")
	}
	//1.创建一个UDP套接字
	conn, err := net.ListenPacket("ip4:udp", srcIP.String())
	if err != nil {
		return fmt.Errorf("ListenPacket失败:%v", err)
	}
	defer conn.Close()
	//2.构造数据包
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
	//3.发送数据包
	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstIP}); err != nil {
		return fmt.Errorf("发送数据包失败:%v", err)
	}
	return nil
}
