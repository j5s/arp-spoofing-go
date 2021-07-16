package main

import (
	"ARPSpoofing/dao/mysql"
	"ARPSpoofing/routers"
	"ARPSpoofing/settings"
	"log"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
)

func main() {
	shell := ishell.NewWithConfig(&readline.Config{
		Prompt: "arp-spoofing > ",
	})

	shell.Println("happy hacking")
	if err := settings.Init(); err != nil {
		log.Println("settings.Init failed,err:", err)
		return
	}
	if err := mysql.Init(); err != nil {
		log.Println("mysql init failed,err:", err)
		return
	}
	routers.Init(shell)
	shell.Run()
}
