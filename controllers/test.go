package controllers

import "github.com/abiosoft/ishell"

func SayHello(c *ishell.Context) {
	c.Println("hello world")
}

func ShowBanner(c *ishell.Context) {
	c.Println("show banner")
}
