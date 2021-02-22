package main

import "github.com/andlabs/ui"

func main() {
	err := ui.Main(setupUI)
	if err != nil {
		panic(err)
	}
}
