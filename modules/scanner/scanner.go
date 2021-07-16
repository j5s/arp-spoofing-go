package modules

import (
	"ARPSpoofing/models"
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
)

type Scanner struct {
	Name    string
	Help    string
	Options *models.Options
	Shell   *ishell.Shell
}

func (s *Scanner) GetName() string {
	return s.Name
}
func (s *Scanner) GetHelp() string {
	return s.Help
}

func (s *Scanner) GetOptions() *models.Options {
	return s.Options
}

func (s *Scanner) CreateShell(c *ishell.Context) {
	s.Shell = ishell.NewWithConfig(&readline.Config{
		Prompt: fmt.Sprintf("arp-spoofing (%s) > ", s.GetName()),
	})
	//展示配置功能
	showCmd := &ishell.Cmd{
		Name: "show",
		Help: "展示信息",
		Func: nil,
	}
	showCmd.AddCmd(&ishell.Cmd{
		Name: "options",
		Help: "展示配置项",
		Func: func(c *ishell.Context) {
			s.Options.Show()
		},
	})
	s.Shell.AddCmd(showCmd)
	//运行模块功能
	s.Shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "运行模块",
		Func: s.Run,
	})
	s.Shell.Run()
}
