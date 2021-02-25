package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/andlabs/ui"
)

func worker(start, end int, iface *net.Interface, hostTable *Table, sender Sender) error {
	hostCh := make(chan *HostItem, 1024)
	//消费者：
	go func() {
		for hostItem := range hostCh {
			if hostItem.MAC.String() == "00:00:00:00:00:00" {
				continue
			}
			if IsContain(hostTable.hModel.Data, *hostItem) == false {
				rows := hostTable.hModel.NumRows(hostTable.model)
				ui.QueueMain(func() {
					hostTable.model.RowInserted(rows)
				})
				hostTable.hModel.Data = append(hostTable.hModel.Data, *hostItem)
			}
		}
	}()

	//消费者+生产者：收包
	go func() {
		if err := sender.Recv(hostCh); err != nil {
			log.Println(err)
		}
	}()

	//生产者：发包
	myIP, _ := getIPv4ByIface(iface)
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
