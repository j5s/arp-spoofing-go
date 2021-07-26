package main

import (
	"ARPSpoofing/dao/redis"
	"ARPSpoofing/debug"
	"ARPSpoofing/logo"
	"ARPSpoofing/routers"
	"ARPSpoofing/settings"
	"log"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
)

func main() {
	shell := ishell.NewWithConfig(&readline.Config{
		Prompt: "阿弥陀佛 > ",
	})
	logo.Show(logo.LogoFile)
	//1.初始化配置
	if err := settings.Init(); err != nil {
		log.Println("settings.Init failed,err:", err)
		return
	}
	//2.连接mysql
	// if err := mysql.Init(); err != nil {
	// 	log.Println("mysql init failed,err:", err)
	// 	return
	// }
	// defer mysql.Close()
	// debug.Println("mysql 数据库连接成功")
	//3.连接redis
	if err := redis.Init(); err != nil {
		log.Println("redis init failed,err:", err)
		return
	}
	defer redis.Close()
	debug.Println("redis 数据库连接成功")

	routers.Init(shell)
	shell.Run()
}
