package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/andlabs/ui"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Config struct {
	Iface      *net.Interface
	MinBox     *ui.Spinbox
	MaxBox     *ui.Spinbox
	ScanMethod ScanMethod
	HideLevel  HideLevel
}
type ScanMethod string

const (
	all ScanMethod = "ALL"
	arp ScanMethod = "ARP request"
	udp ScanMethod = "UDP+ICMP"
)

var methods []ScanMethod = []ScanMethod{
	all,
	arp,
	udp,
}

type HideLevel string

const (
	positive        = "主动扫描"
	passive         = "被动扫描"
	positivePassive = "主动被动结合"
)

var hideLevels []HideLevel = []HideLevel{
	positivePassive,
	positive,
	passive,
}

type Sender interface {
	Send(dstIP net.IP) error
	Recv(out chan *HostItem) error
}
type ARPSender struct {
	ethDstMAC net.HardwareAddr
	ethSrcMAC net.HardwareAddr
	arpDstMAC net.HardwareAddr
	arpSrcMAC net.HardwareAddr
	arpSrcIP  net.IP
	handle    *pcap.Handle
}

func newARPSender(iface *net.Interface) (Sender, error) {
	log.Println("new ARP Sender")
	handle, err := pcap.OpenLive(iface.Name, 65535, true, 1000)
	if err != nil {
		return nil, fmt.Errorf("OpenLive Error:%v", err)
	}
	myIP, err := getIPv4ByIface(iface)
	if err != nil {
		return nil, fmt.Errorf("IPv4 error:%v", err)
	}
	ethDstMAC, err := net.ParseMAC("ff:ff:ff:ff:ff:ff")
	if err != nil {
		return nil, err
	}
	arpDstMAC, err := net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		return nil, err
	}
	return &ARPSender{
		ethSrcMAC: iface.HardwareAddr,
		ethDstMAC: ethDstMAC,
		arpSrcMAC: iface.HardwareAddr,
		arpDstMAC: arpDstMAC,
		arpSrcIP:  myIP,
		handle:    handle,
	}, nil
}
func (this *ARPSender) Send(dstIP net.IP) error {
	return SendARP(this.handle, this.ethDstMAC, this.ethSrcMAC, this.arpDstMAC, this.arpSrcMAC, dstIP, this.arpSrcIP, layers.ARPRequest)
}
func (this *ARPSender) Recv(out chan *HostItem) error {
	ReceiveARP(this.handle, this.ethSrcMAC, out, true)
	return nil
}

type UDPSender struct {
	srcIP   net.IP
	dstPort int
	srcPort int
	dataLen int
	iface   *net.Interface
}

func newUDPSender(iface *net.Interface) (Sender, error) {
	myIP, err := getIPv4ByIface(iface)
	if err != nil {
		return nil, err
	}
	return &UDPSender{
		srcIP:   myIP,
		dstPort: 65534,
		srcPort: rand.Intn(65535),
		dataLen: 1000,
		iface:   iface,
	}, nil
}
func (this *UDPSender) Send(dstIP net.IP) error {
	return SendUDPv2(this.srcIP.String(), dstIP.String(), this.srcPort, this.dstPort, this.dataLen)
}
func (this *UDPSender) Recv(out chan *HostItem) error {
	handle, err := pcap.OpenLive(this.iface.Name, 65535, true, 1000)
	if err != nil {
		return fmt.Errorf("OpenLive Error:%v", err)
	}
	RecvICMP(handle, out)
	return nil
}

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
