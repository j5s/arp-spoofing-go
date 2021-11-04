package controllers

import (
	"ARPSpoofing/dao/redis"
	"ARPSpoofing/logic"
	"ARPSpoofing/settings"
	"ARPSpoofing/vars"
	"context"
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
)

//MiddleAttackHandler 中间人攻击
func MiddleAttackHandler(c *ishell.Context) {
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
	noAttackedList := make([]string, 0, len(ipList))
	for i := range ipList {
		if ipList[i] == gateway {
			continue
		}
		if _, exist := vars.MiddleAttackCancelMap[ipList[i]]; true == exist {
			continue
		}
		noAttackedList = append(noAttackedList, ipList[i])
	}
	if len(noAttackedList) == 0 {
		fmt.Println("暂时没有可以攻击的主机")
		return
	}
	targetIndex := c.MultiChoice(noAttackedList, "which host do you want to attack?")
	c.Println(color.YellowString(fmt.Sprintf("[*] 成为%s 和%s之间的中间人", gateway, noAttackedList[targetIndex])))
	ctx, cancel := context.WithCancel(context.Background())
	if err := logic.MiddleAttack(ctx, gateway, noAttackedList[targetIndex], gateway); err != nil {
		c.Println(color.RedString(fmt.Sprintf("[*] 中间人攻击协程启动失败,logic.MiddleAttack failed,err:%v", err)))
		_ = cancel
		return
	}
	c.Println(color.GreenString("[*] 中间人攻击协程启动成功"))
	c.Println(color.GreenString("[*] show middle_attacked 查看正在进行的中间人攻击"))
	c.Println(color.GreenString("[*] stop middle_attacked 停止中间人攻击"))
	//3.维护一个中间人攻击和退出函数的映射
	vars.MiddleAttackCancelMap[noAttackedList[targetIndex]] = cancel
}

//StopMiddleAttackHandler 停止中间人攻击
func StopMiddleAttackHandler(c *ishell.Context) {
	hosts := make([]string, 0)
	for key := range vars.MiddleAttackCancelMap {
		hosts = append(hosts, key)
	}
	if len(hosts) == 0 {
		c.Println("暂时没有正在受攻击的主机")
		return
	}
	choice := c.MultiChoice(hosts, "Please select host.")
	c.Println("Stop attacking:", hosts[choice])
	vars.MiddleAttackCancelMap[hosts[choice]]()
	delete(vars.MiddleAttackCancelMap, hosts[choice])
}
