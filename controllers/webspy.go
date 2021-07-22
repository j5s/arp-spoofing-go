package controllers

import (
	"ARPSpoofing/pkg/assembly"
	"ARPSpoofing/pkg/server"
	"context"

	"github.com/abiosoft/ishell"
	"github.com/spf13/viper"
)

var webSpyCancelFunc context.CancelFunc = nil

//WebSpyHandler 嗅探Web 请求
func WebSpyHandler(c *ishell.Context) {
	if webSpyCancelFunc != nil {
		c.Println("[*] webspy 已启动")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	webSpyCancelFunc = cancel
	//1.启动服务器 用于向用户显示数据
	go server.Run(ctx)
	c.Printf("[*] 启动服务器成功:http://127.0.0.1:%d/index\n", viper.GetInt("webspy.port"))
	//2.启动监听程序
	go assembly.Run(ctx)
	c.Println("[*] 嗅探程序启动成功")
}

//StopWebSpyHandler 停止嗅探
func StopWebSpyHandler(c *ishell.Context) {
	if webSpyCancelFunc != nil {
		webSpyCancelFunc()
		webSpyCancelFunc = nil
	}
}
