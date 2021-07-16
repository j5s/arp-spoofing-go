package routers

import (
	"ARPSpoofing/controllers"
	"ARPSpoofing/settings"

	"github.com/abiosoft/ishell"
)

//Init 初始化路由
func Init(shell *ishell.Shell) {
	//未找到命令时
	shell.NotFound(controllers.NotFoundHandler)
	//show 展示功能
	showCmd := &ishell.Cmd{
		Name: "show",
		Help: "展示信息",
		Func: nil,
	}
	showCmd.AddCmd(&ishell.Cmd{
		Name: "options",
		Help: "展示配置项",
		Func: controllers.ShowOptionsHandler,
	})
	shell.AddCmd(showCmd)

	//set 设置配置功能
	setCmd := &ishell.Cmd{
		Name: "set",
		Help: "配置参数",
		Func: nil,
	}
	for _, item := range settings.Options.Content {
		setCmd.AddCmd(&ishell.Cmd{
			Name: item.Name,
			Help: item.Description,
			Func: controllers.SetOptionHandler,
		})
	}
	shell.AddCmd(setCmd)

	//scan 扫描功能
	shell.AddCmd(&ishell.Cmd{
		Name: "scan",
		Help: "扫描内网中存活的主机",
		Func: controllers.ScanHandler,
	})
}
