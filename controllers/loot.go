package controllers

import (
	"ARPSpoofing/logic"

	"github.com/abiosoft/ishell"
)

//ShowLootHandler 展示所有的战利品
func ShowLootHandler(c *ishell.Context) {
	if err := logic.ShowLoot(); err != nil {
		c.Println("logic.ShowHosts() failed,err:", err)
	}
}

//ClearLootHandler 清除所有战利品
func ClearLootHandler(c *ishell.Context) {
	if err := logic.ClearLoot(); err != nil {
		c.Println("logic.ClearLoot() failed,err:", err)
	}
	c.Println("[*] clear loot success")
}
