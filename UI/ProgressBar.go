package UI

import (
	"ARPSpoofing/Socket"

	"github.com/andlabs/ui"
)

func makeProgressDiv(config *Socket.Config) (*ui.Box, *ui.Button) {
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
