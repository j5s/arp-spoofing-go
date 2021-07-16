package controllers

import (
	"os/exec"
	"strings"

	"github.com/abiosoft/ishell"
)

//NotFoundHandler 未找到命令时处理函数
func NotFoundHandler(c *ishell.Context) {
	cmd := exec.Command("/bin/bash", "-c", strings.Join(c.Args, " "))
	result, err := cmd.CombinedOutput()
	if err != nil {
		c.Println(err)
		return
	}
	c.Printf(string(result))
}
