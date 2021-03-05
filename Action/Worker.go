package Action

import (
	"ARPSpoofing/Socket"
	"ARPSpoofing/Utils"

	"fmt"
	"log"
	"net"
	"strings"

	"github.com/andlabs/ui"
)

func Worker(start, end int, iface *net.Interface, hostTable *Table, sender Socket.Sender, hideLevel Socket.HideLevel) error {
	hostCh := make(chan *Socket.HostItem, 1024)
	//消费者：
	go func() {
		for hostItem := range hostCh {
			if hostItem.MAC.String() == "00:00:00:00:00:00" {
				continue
			}
			if Socket.IsContain(hostTable.HModel.Data, *hostItem) == false {
				rows := hostTable.HModel.NumRows(hostTable.Model)
				ui.QueueMain(func() {
					hostTable.Model.RowInserted(rows)
				})
				hostTable.HModel.Data = append(hostTable.HModel.Data, *hostItem)
			}
		}
	}()

	//消费者+生产者：收包
	go func() {
		if err := sender.Recv(hostCh); err != nil {
			log.Println(err)
		}
	}()
	if hideLevel == Socket.Passive {
		log.Println("被动扫描不发包，只收包")
		return nil
	}
	//生产者：发包
	myIP, _ := Utils.GetIPv4ByIface(iface)
	myIParr := strings.Split(myIP.String(), ".")
	for i := start; i <= end; i++ {
		dstIP := net.ParseIP(fmt.Sprintf("%s.%s.%s.%d", myIParr[0], myIParr[1], myIParr[2], i))
		err := sender.Send(dstIP)
		if err != nil {
			log.Println("发送数据包失败:", err)
			continue
		}
	}
	return nil
}
