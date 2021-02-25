package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

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
			out <- newHostItem(net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))
			out <- newHostItem(net.IP(arp.DstProtAddress), net.HardwareAddr(arp.DstHwAddress))
		} else {
			if arp.Operation != layers.ARPReply {
				continue
			}
			if !bytes.Equal([]byte(myMAC), arp.DstHwAddress) {
				continue
			}
			out <- newHostItem(net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))
		}
	}
	log.Println("ARPRecive quit")
}

//获取分配给该网卡的内网ipv4地址
func getIPv4ByIface(iface *net.Interface) (net.IP, error) {
	addrs, err := iface.Addrs() //ipv6/mask,ipv4/mask
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ipmask, ok := addr.(*net.IPNet)
		if ok {
			ipv4 := ipmask.IP.To4()
			if ipv4 != nil {
				return ipv4, nil
			}
		}
	}
	return nil, errors.New("don't have ipv4 address")
}

func CutOff(timeCh <-chan time.Time, hostTable *Table, iface net.Interface) error {
	handle, err := pcap.OpenLive(iface.Name, 65536, true, 1000)
	if err != nil {
		return err
	}
	for {
		<-timeCh
		for _, host := range hostTable.hModel.Data {
			if host.IsCutOff == false {
				continue
			}
			var dstIP net.IP
			var dstMAC net.HardwareAddr
			var srcIP net.IP
			var srcMAC net.HardwareAddr = net.HardwareAddr{0x66, 0x66, 0x66, 0x66, 0x66, 0x66}
			gIndex := hostTable.hModel.GatewayIndex
			var gIP net.IP = net.IP(hostTable.hModel.Data[gIndex].IP)
			var gMAC net.HardwareAddr = net.HardwareAddr(hostTable.hModel.Data[gIndex].MAC)
			var op uint16
			if host.Spooling == "Host" {
				dstIP = net.IP(host.IP)
				dstMAC = net.HardwareAddr(host.MAC)
				srcIP = gIP
			} else {
				dstIP = gIP
				dstMAC = gMAC
				srcIP = net.IP(host.IP)
			}
			if host.PacketType == "Request" {
				op = layers.ARPRequest
			} else {
				op = layers.ARPReply
			}
			err := SendARP(handle, dstMAC, srcMAC, dstMAC, srcMAC, dstIP, srcIP, op)
			if err != nil {
				log.Println("sent arp request failed,because ", err)
			} else {
				log.Printf("ethDstMAC:%s,ethSrcMAC:%s,arpDstMAC:%s,arpSrcIP:%s,arpDstIP:%s,arpSrcIP:%s\n",
					dstMAC.String(), srcMAC.String(),
					dstMAC.String(), srcMAC.String(),
					dstIP.String(), srcIP.String())
			}
		}
	}
}
