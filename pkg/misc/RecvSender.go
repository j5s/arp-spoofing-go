package misc

import (
	"ARPSpoofing/models"
	"ARPSpoofing/pkg/arp"
	"ARPSpoofing/pkg/udp"
	"context"
	"net"
)

//RecvSender 混合发送器
type RecvSender struct {
	udpSender models.RecvSender
	arpSender models.RecvSender
}

//NewRecvSender 构造函数
func NewRecvSender(iface *net.Interface) (models.RecvSender, error) {
	udpSender, err := udp.NewRecvSender(iface)
	if err != nil {
		return nil, err
	}
	arpSender, err := arp.NewRecvSender(iface)
	if err != nil {
		return nil, err
	}
	return &RecvSender{
		udpSender: udpSender,
		arpSender: arpSender,
	}, nil
}

//Send 发送数据包
func (s *RecvSender) Send(dstIP net.IP) error {
	if err := s.arpSender.Send(dstIP); err != nil {
		return err
	}
	if err := s.udpSender.Send(dstIP); err != nil {
		return err
	}
	return nil
}

//Recv 接收数据包
func (s *RecvSender) Recv(ctx context.Context) <-chan *models.Host {
	out := make(chan *models.Host, 2048)
	//1.启动两个接收协程 并返回管道
	udpOut := s.udpSender.Recv(ctx)
	arpOut := s.arpSender.Recv(ctx)
	//2.不断从管道中读数据
	go func() {
		for {
			select {
			case host := <-udpOut:
				out <- host
			case host := <-arpOut:
				out <- host
			}
		}
	}()
	return out
}
