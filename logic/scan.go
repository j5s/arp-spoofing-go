package logic

import (
	"ARPSpoofing/dao/redis"
	"ARPSpoofing/debug"
	"ARPSpoofing/models"
	"ARPSpoofing/pkg/arp"
	"ARPSpoofing/pkg/misc"
	"ARPSpoofing/pkg/routine"
	"ARPSpoofing/pkg/udp"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/abiosoft/ishell"
	manuf "github.com/timest/gomanuf"
)

//Scan ARP扫描
func Scan(c *ishell.Context, ipList []net.IP, ifname string, method string) error {
	sender, err := NewRecvSender(method, ifname)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	// 消费者：向数据库中存储主机信息
	go func(ctx context.Context) {
		debug.Println("启动了一个存储数据的协程:", routine.GetGID())
		defer debug.Println("存储数据的协程退出:", routine.GetGID())
		for host := range sender.Recv(ctx) {
			//查询mac地址对应的厂商信息
			host.MACInfo = manuf.Search(host.MAC)
			err = redis.NewHosts().Add(host)
			if err != nil {
				log.Println("redis add Host failed,err:", err)
				return
			}
			select {
			case <-ctx.Done(): // 等待上级通知退出
				return
			default:
			}
		}
	}(ctx)
	// 生产者：发包
	c.ProgressBar().Start()
	total := len(ipList)
	for index, dstIP := range ipList {
		if err := sender.Send(dstIP); err != nil {
			log.Println("发送数据包失败:", err)
		}
		c.ProgressBar().Suffix(fmt.Sprint(" ", index*100/total, "%"))
		c.ProgressBar().Progress(index * 100 / total)
		time.Sleep(time.Millisecond * 30)
	}
	c.ProgressBar().Stop()
	cancel()
	debug.Println("通知子协程退出")
	return nil
}

//NewRecvSender 简单工厂模式
func NewRecvSender(scanMethod string, ifname string) (models.RecvSender, error) {
	switch scanMethod {
	case "arp":
		return arp.NewRecvSender(ifname)
	case "udp":
		return udp.NewRecvSender(ifname)
	case "all":
		return misc.NewRecvSender(ifname)
	default:
		return nil, errors.New("error scanMethod")
	}
}
