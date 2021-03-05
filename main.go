package main

import (
	"ARPSpoofing/UI"

	"github.com/andlabs/ui"
)

func main() {
	err := ui.Main(UI.SetupUI)
	if err != nil {
		panic(err)
	}
}
