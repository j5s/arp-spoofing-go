package main

import (
	"log"

	"github.com/andlabs/ui"
)

func setupUI() {
	window := ui.NewWindow("ARPSpoofing Author:XDUer RickShang", 800, 600, false)
	window.SetMargined(true)
	BodyDiv := ui.NewVerticalBox()
	BodyDiv.SetPadded(true)
	//配置面板
	configForm, config := makeConfigForm(window)
	BodyDiv.Append(configForm, true)
	//输出面板
	hostTable := makeHostTable()
	BodyDiv.Append(hostTable.table, false)
	//进度条
	progressDiv, startBtn := makeProgressDiv(config)
	//ARP Spoofing
	go func() {
		timeTicker := hostTable.hModel.timeTicker
		CutOff(timeTicker, hostTable, *config.Iface)
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
			if err := worker(start, end, config.Iface, hostTable, sender); err != nil {
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

func makeProgressDiv(config *Config) (*ui.Box, *ui.Button) {
	progressDiv := ui.NewVerticalBox()
	progressBar := ui.NewProgressBar()
	progressDiv.SetPadded(true)
	bottomDiv := ui.NewHorizontalBox()
	bottomDiv.SetPadded(true)
	currentLabel := ui.NewLabel("Current Scanning...")
	startBtn := ui.NewButton("Scan")

	bottomDiv.Append(currentLabel, true)
	bottomDiv.Append(startBtn, false)

	progressDiv.Append(progressBar, false)
	progressDiv.Append(bottomDiv, false)

	return progressDiv, startBtn
}

func makeHostTable() *Table {
	hModel := newModelHandler()
	hostModel := ui.NewTableModel(hModel)
	hostTable := ui.NewTable(&ui.TableParams{
		Model:                         hostModel,
		RowBackgroundColorModelColumn: 3, //-1 默认颜色
	})
	hostTable.AppendTextColumn("IP", 0, 0, nil)
	hostTable.AppendTextColumn("MAC", 1, 0, nil)
	hostTable.AppendTextColumn("MACInfo", 2, 0, nil)
	hostTable.AppendButtonColumn("IsGateway", 3, ui.TableModelColumnAlwaysEditable)
	hostTable.AppendButtonColumn("Spooling", 4, ui.TableModelColumnAlwaysEditable)
	hostTable.AppendButtonColumn("PacketType", 5, ui.TableModelColumnAlwaysEditable)
	hostTable.AppendButtonColumn("Cut Off", 6, ui.TableModelColumnAlwaysEditable)

	return &Table{
		table:  hostTable,
		model:  hostModel,
		hModel: hModel,
	}
}
