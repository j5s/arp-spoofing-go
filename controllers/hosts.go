package controllers

import (
	"ARPSpoofing/logic"

	"github.com/abiosoft/ishell"
)

//ShowHostsHandler 展示所有主机
func ShowHostsHandler(c *ishell.Context) {
	err := logic.ShowHosts()
	if err != nil {
		c.Println("logic.ShowHosts failed,err:", err)
	}
}

//ClearHostsHandler 清空所有主机
func ClearHostsHandler(c *ishell.Context) {
	err := logic.ClearHosts()
	if err != nil {
		c.Println("logic.ClearHosts failed,err:", err)
	}
	c.Println("清空所有主机成功")
}
