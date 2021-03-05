package Socket

import (
	"errors"
	"net"

	"github.com/andlabs/ui"
)

type Config struct {
	Iface      *net.Interface
	MinBox     *ui.Spinbox
	MaxBox     *ui.Spinbox
	ScanMethod ScanMethod
	HideLevel  HideLevel
}

//简单工厂模式
func (c *Config) NewSender() (Sender, error) {
	switch c.ScanMethod {
	case arp:
		return newARPSender(c.Iface)
	case udp:
		return newUDPSender(c.Iface)
	case all:
		return newSuperSender(c.Iface)
	default:
		return nil, errors.New("error scanMethod")
	}
}
