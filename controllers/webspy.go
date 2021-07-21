package controllers

import (
	"ARPSpoofing/pkg/assembly"
	"ARPSpoofing/pkg/server"
	"fmt"

	"github.com/abiosoft/ishell"
)

//WebSpyHandler 嗅探Web 请求
func WebSpyHandler(c *ishell.Context) {
	//1.启动服务器 用于向用户显示数据
	server.Run()
	//2.启动监听程序
	go assembly.Run()
	fmt.Println("监听程序启动成功")
}
