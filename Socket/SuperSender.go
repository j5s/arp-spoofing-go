package Socket

import (
	"log"
	"net"
)

type SuperSender struct {
	udpSender Sender
	arpSender Sender
}

func newSuperSender(iface *net.Interface) (Sender, error) {
	udpSender, err := newUDPSender(iface)
	if err != nil {
		return nil, err
	}
	arpSender, err := newARPSender(iface)
	if err != nil {
		return nil, err
	}
	return &SuperSender{
		udpSender: udpSender,
		arpSender: arpSender,
	}, nil
}
func (this *SuperSender) Send(dstIP net.IP) error {
	// log.Println("SuperSender Send", dstIP.String())
	// log.Println("send arp")
	if err := this.arpSender.Send(dstIP); err != nil {
		return err
	}
	// log.Println("send udp")
	if err := this.udpSender.Send(dstIP); err != nil {
		return err
	}
	return nil
}
func (this *SuperSender) Recv(out chan *HostItem) error {

	go func() {
		log.Println("recv arp start")
		err := this.arpSender.Recv(out)
		if err != nil {
			log.Println(err)
		}
		log.Println("recv arp done")
	}()
	go func() {
		log.Println("recv udp start")
		err := this.udpSender.Recv(out)
		if err != nil {
			log.Println(err)
		}
		log.Println("recv udp done")
	}()
	return nil
}
