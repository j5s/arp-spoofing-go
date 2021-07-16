package models

import (
	"net"

	"github.com/abiosoft/ishell"
)

//ModuleItem 模块
type ModuleItem interface {
	GetName() string
	GetHelp() string
	GetOptions() Options
	CreateShell(c *ishell.Context)
	Run(c *ishell.Context)
}

//Sender 凡是能发送包，接收包都是Sender
type Sender interface {
	Send(dstIP net.IP) error
	Recv(out chan *Host) error
}
