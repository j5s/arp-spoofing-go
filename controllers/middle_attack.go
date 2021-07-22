package controllers

import (
	"ARPSpoofing/dao/redis"
	"ARPSpoofing/logic"
	"ARPSpoofing/settings"
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
	targetIndex := c.MultiChoice(ipList, "which host do you want to attack?")
	fmt.Printf("\r成为%s 和%s之间的中间人\n", gateway, ipList[targetIndex])
	if err := logic.MiddleAttack(gateway, ipList[targetIndex], gateway); err != nil {
		log.Println("logic.MiddleAttack failed,err:", err)
		return
	}
}
