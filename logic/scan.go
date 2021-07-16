package logic

import (
	"ARPSpoofing/dao/mysql"
	"ARPSpoofing/models"
	"ARPSpoofing/socket"
	"log"
	"net"
)

//Scan ARP扫描
func Scan(ipList []net.IP, iface *net.Interface, method string) error {
	sender, err := socket.NewSender(method, iface)
	if err != nil {
		log.Println(err)
		return err
	}
	hostCh := make(chan *models.Host, 1024)
	//消费者：向数据库中存储主机信息
	go func() {
		for host := range hostCh {
			err := mysql.AddHost(host)
			if err != nil {
				log.Println("mysql.AddHost failed,err:", err)
				return
			}
		}
	}()
	//消费者+生产者：收包
	go func() {
		if err := sender.Recv(hostCh); err != nil {
			log.Println(err)
			return
		}
	}()
	//生产者：发包
	for dstIP := range ipList {
		if err := sender.Send(dstIP); err != nil {
			log.Println("发送数据包失败:", err)
		}
	}
	return nil
}
