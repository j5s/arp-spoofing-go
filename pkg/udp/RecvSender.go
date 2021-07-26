package udp

import (
	"ARPSpoofing/models"
	"ARPSpoofing/pkg/icmp"
	"ARPSpoofing/utils"
	"context"
	"fmt"
	"math/rand"
	"net"

	"github.com/google/gopacket/pcap"
)

//RecvSender 收发器
type RecvSender struct {
	srcIP   net.IP
	srcPort int
	dstPort int
	dataLen int
	handle  *pcap.Handle
}

//NewRecvSender 构造函数
func NewRecvSender(ifname string) (models.RecvSender, error) {
	handle, err := pcap.OpenLive(ifname, 65535, true, 1000)
	if err != nil {
		return nil, fmt.Errorf("OpenLive Error:%v", err)
	}
	myIP, err := utils.GetIPv4(ifname)
	if err != nil {
		return nil, err
	}
	return &RecvSender{
		srcIP:   myIP,
		dstPort: 65534,
		srcPort: rand.Intn(65535),
		dataLen: 1000,
		handle:  handle,
	}, nil
}

//Send 发送数据包
func (r *RecvSender) Send(dstIP net.IP) error {
	return SendUDPv2(r.srcIP.String(), dstIP.String(), r.srcPort, r.dstPort, r.dataLen)
}

//Recv 接收数据包
func (r *RecvSender) Recv(ctx context.Context) <-chan *models.Host {
	return icmp.RecvICMP(ctx, r.handle)
}
