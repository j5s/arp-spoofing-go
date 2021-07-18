package controllers

import (
	"ARPSpoofing/logic"
	"ARPSpoofing/vars"

	"github.com/abiosoft/ishell"
)

//SniffHandler 嗅探功能
func SniffHandler(c *ishell.Context) {
	//1.效验
	if vars.SniffCancelFunc != nil {
		c.Println("已经有一个敏感报文嗅探器在后台工作,请先退出它")
		return
	}
	//2.业务逻辑
	if err := logic.Sniff(); err != nil {
		c.Println("logic.sniff failed,err:", err)
		return
	}
	c.Println("[*] 启动敏感报文嗅探器成功")
}

//StopSniffHandler 停止嗅探器
func StopSniffHandler(c *ishell.Context) {
	if vars.SniffCancelFunc == nil {
		c.Println("[*] 嗅探器未启动")
		return
	}
	vars.SniffCancelFunc()
	c.Println("[*] 停止嗅探器成功")
}

//CheckSniffHandler 检查嗅探器的工作状态
func CheckSniffHandler(c *ishell.Context) {
	statusMap := map[bool]string{
		true:  "sniffer is running    [ON]",
		false: "sniffer is stopped    [OFF]",
	}
	c.Println(statusMap[vars.SniffCancelFunc != nil])
}
