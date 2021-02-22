package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/andlabs/ui"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func worker(start, end int, iface *net.Interface, hostTable *Table) error {
	handle, err := pcap.OpenLive(iface.Name, 65536, true, 1000)
	if err != nil {
		return err
	}
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
		for {
			ReceiveARP(handle, iface, hostCh, true)
		}
	}()
	//生产者：发包
	myIP, _ := getIPv4ByIface(iface)
	ipSlice := strings.Split(myIP.String(), ".")
	nums := []int{}
	for i := 0; i < 3; i++ {
		i, err := strconv.Atoi(ipSlice[i])
		if err != nil {
			return err
		}
		nums = append(nums, i)
	}
	fmt.Println("************", nums)
	ethDstMAC := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	ethSrcMAC := iface.HardwareAddr
	arpDstMAC := net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	arpSrcMAC := ethSrcMAC
	arpSrcIP, _ := getIPv4ByIface(iface)
	for i := start; i <= end; i++ {
		_arpDstIP := net.IPv4(byte(nums[0]), byte(nums[1]), byte(nums[2]), byte(i))
		arpDstIP := host2net(_arpDstIP)
		err := SendARP(handle, ethDstMAC, ethSrcMAC, arpDstMAC, arpSrcMAC, arpDstIP, arpSrcIP, layers.ARPRequest)
		if err != nil {
			log.Println("sent arp request failed,because ", err)
		} else {
			log.Printf("ethDstMAC:%s,ethSrcMAC:%s,arpDstMAC:%s,arpSrcIP:%s,arpDstIP:%s,arpSrcIP:%s\n",
				ethDstMAC.String(), ethSrcMAC.String(),
				arpDstMAC.String(), arpSrcMAC.String(),
				_arpDstIP, arpSrcIP.String())
		}
	}
	//给1000s的收包时间
	time.Sleep(10 * time.Second)
	return nil
}
