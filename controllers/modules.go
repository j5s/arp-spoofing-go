package controllers

import (
	"ARPSpoofing/vars"

	"github.com/abiosoft/ishell"
)

func ShowModulesHandler(c *ishell.Context) {
	vars.Mod.Show()
}
