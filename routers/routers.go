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
	//host 主机管理
	hostCmd := &ishell.Cmd{
		Name: "hosts",
		Help: "主机管理功能",
		Func: nil,
	}
	hostCmd.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "展示所有主机",
		Func: controllers.ShowHostsHandler,
	})
	hostCmd.AddCmd(&ishell.Cmd{
		Name: "clear",
		Help: "清空所有主机",
		Func: controllers.ClearHostsHandler,
	})
	shell.AddCmd(hostCmd)
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
	//cut 断网功能
	cutCmd := &ishell.Cmd{
		Name: "cut",
		Help: "通过ARP欺骗切断局域网内某台主机的网络",
		Func: nil,
	}
	cutCmd.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "开始攻击",
		Func: controllers.CutHandler,
	})
	cutCmd.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "停止攻击",
		Func: controllers.StopCutHandler,
	})
	shell.AddCmd(cutCmd)
	//sniff 嗅探功能
	sniffCmd := &ishell.Cmd{
		Name: "sniff",
		Help: "嗅探用户名和密码",
		Func: nil,
	}
	sniffCmd.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "启动敏感报文嗅探器",
		Func: controllers.SniffHandler,
	})
	sniffCmd.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "停止敏感报文嗅探器",
		Func: controllers.StopSniffHandler,
	})
	sniffCmd.AddCmd(&ishell.Cmd{
		Name: "status",
		Help: "查看敏感报文嗅探器的状态",
		Func: controllers.CheckSniffHandler,
	})
	shell.AddCmd(sniffCmd)
	//loot 敏感信息
	shell.AddCmd(&ishell.Cmd{})
}
