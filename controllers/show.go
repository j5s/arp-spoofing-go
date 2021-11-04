package controllers

import (
	"ARPSpoofing/pkg/table"
	"ARPSpoofing/settings"
	"ARPSpoofing/vars"
	"fmt"
	"log"

	"github.com/abiosoft/ishell"
	"github.com/google/gopacket/pcap"
)

//SetOptionHandler 设置选项
func SetOptionHandler(c *ishell.Context) {
	//1.效验参数
	if len(c.Args) != 1 {
		c.Println("[*] set name value")
		c.Println(c.Cmd.Help)
		return
	}
	//2.业务逻辑
	settings.Options.Set(c.Cmd.Name, c.Args[0])
}

//ShowOptionsHandler 展示所有配置项
func ShowOptionsHandler(c *ishell.Context) {
	settings.Options.Show()
}

//ShowIfnamesHandler 展示所有网卡名
func ShowIfnamesHandler(c *ishell.Context) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Println("pcap.FindAllDevs failed,err:", err)
		return
	}
	for i := range devices {
		if len(devices[i].Addresses) == 0 {
			continue
		}
		fmt.Printf("Interface Name:%s\n", devices[i].Name)
		for _, addr := range devices[i].Addresses {
			fmt.Println("|- IP address:", addr.IP)
			fmt.Println("|- Subnet mask:", addr.Netmask)
		}
		fmt.Println()
	}
}

//ShowCuttedHandler 显示所有被切断的主机
func ShowCuttedHandler(c *ishell.Context) {
	if len(vars.HostCancelMap) == 0 {
		c.Println("暂无被切断的主机")
		return
	}
	headers := []string{"Cutted Hosts"}
	data := make([][]string, 0, len(vars.HostCancelMap))
	for key := range vars.HostCancelMap {
		data = append(data, []string{key})
	}
	table.Show(headers, data)
}

//ShowMiddleAttackedHandler 显示所有被中间人攻击的主机
func ShowMiddleAttackedHandler(c *ishell.Context) {
	if len(vars.MiddleAttackCancelMap) == 0 {
		c.Println("暂无被中间人攻击的主机")
		return
	}
	headers := []string{"MiddleAttacked Hosts"}
	data := make([][]string, 0, len(vars.MiddleAttackCancelMap))
	for key := range vars.MiddleAttackCancelMap {
		data = append(data, []string{key})
	}
	table.Show(headers, data)
}
