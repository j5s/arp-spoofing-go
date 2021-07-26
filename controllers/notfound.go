package controllers

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/abiosoft/ishell"
)

//NotFoundHandler 未找到命令时处理函数
func NotFoundHandler(c *ishell.Context) {
	var cmd *exec.Cmd
	input := strings.Join(c.Args, " ")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", input) //windows
	} else {
		cmd = exec.Command("/bin/bash", "-c", input)
	}
	result, err := cmd.CombinedOutput()
	if err != nil {
		c.Println(err)
		return
	}
	c.Printf(string(result))
}
