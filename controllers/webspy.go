package controllers

import (
	"ARPSpoofing/pkg/server"

	"github.com/abiosoft/ishell"
)

//WebSpyHandler 嗅探Web 请求
func WebSpyHandler(c *ishell.Context) {
	//1.启动服务器 用于向用户显示数据
	server.Run()
	//2.启动监听程序
}
