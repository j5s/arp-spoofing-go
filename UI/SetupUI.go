package UI

import (
	"ARPSpoofing/Action"
	"log"

	"github.com/andlabs/ui"
)

func SetupUI() {
	window := ui.NewWindow("ARPSpoofing Author:XDUer RickShang", 800, 600, false)
	window.SetMargined(true)
	BodyDiv := ui.NewVerticalBox()
	BodyDiv.SetPadded(true)
	//配置面板
	configForm, config := Action.MakeConfigForm(window)
	BodyDiv.Append(configForm, true)
	//输出面板
	hostTable := Action.MakeHostTable()
	BodyDiv.Append(hostTable.Table, false)
	//进度条
	progressDiv, startBtn := makeProgressDiv(config)
	//ARP Spoofing
	go func() {
		timeTicker := hostTable.HModel.TimeTicker
		Action.CutOff(timeTicker, hostTable, *config.Iface)
	}()
	//开始扫描
	startBtn.OnClicked(func(self *ui.Button) {
		end := config.MaxBox.Value()
		start := config.MinBox.Value()

		sender, err := config.NewSender()
		if err != nil {
			panic(err)
		}

		go func() {
			if err := Action.Worker(start, end, config.Iface, hostTable, sender, config.HideLevel); err != nil {
				log.Println(err)
			}
		}()
	})
	BodyDiv.Append(progressDiv, false)
	window.SetChild(BodyDiv)
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	window.Show()
}
