package Socket

import (
	"net"

	manuf "github.com/timest/gomanuf"
)

type HostItem struct {
	IP         net.IP
	MAC        net.HardwareAddr
	MACInfo    string
	Spooling   string
	IsCutOff   bool
	PacketType string
}

func NewHostItem(IP net.IP, MAC net.HardwareAddr) *HostItem {
	return &HostItem{
		IP:         IP,
		MAC:        MAC,
		MACInfo:    manuf.Search(MAC.String()),
		Spooling:   "Host",
		IsCutOff:   false,
		PacketType: "Reply",
	}
}

func IsContain(hosts []HostItem, host HostItem) bool {
	for _, eachHost := range hosts {
		if eachHost.IP.String() == host.IP.String() && eachHost.MAC.String() == host.MAC.String() {
			return true
		}
	}
	return false
}
