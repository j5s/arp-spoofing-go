package models

import (
	"context"
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

//RecvSender 凡是能发送包，接收包都是RevSender
type RecvSender interface {
	Send(dstIP net.IP) error
	Recv(ctx context.Context) <-chan *Host
}
