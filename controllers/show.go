package controllers

import (
	"ARPSpoofing/settings"

	"github.com/abiosoft/ishell"
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
