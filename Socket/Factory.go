package socket

import (
	"errors"
	"net"
)

//NewSender 简单工厂模式
func NewSender(scanMethod string, iface *net.Interface) (Sender, error) {
	switch scanMethod {
	case "arp":
		return newARPSender(iface)
	case "udp":
		return newUDPSender(iface)
	case "all":
		return newSuperSender(iface)
	default:
		return nil, errors.New("error scanMethod")
	}
}
