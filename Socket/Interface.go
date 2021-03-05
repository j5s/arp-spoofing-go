package Socket

import (
	"net"
)

type ScanMethod string

const (
	all ScanMethod = "ALL"
	arp ScanMethod = "ARP request"
	udp ScanMethod = "UDP+ICMP"
)

var Methods []ScanMethod = []ScanMethod{
	all,
	arp,
	udp,
}

type Sender interface {
	Send(dstIP net.IP) error
	Recv(out chan *HostItem) error
}

type HideLevel string

const (
	Positive = "主动扫描"
	Passive  = "被动扫描"
)

var HideLevels []HideLevel = []HideLevel{
	Positive,
	Passive,
}
