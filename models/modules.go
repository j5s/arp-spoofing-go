package models

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

//Modules 所有模块
type Modules []ModuleItem

//Add 添加模块到模块表
func (m *Modules) Add(item ModuleItem) {
	*m = append(*m, item)
}

//Show 展示所有模块
func (m Modules) Show() {
	lines := make([][]string, 0)
	for index, item := range m {
		lines = append(lines, []string{
			fmt.Sprintf("%v", index),
			item.GetName(),
			item.GetHelp(),
		})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "NAME", "DESCRIPTION"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.AppendBulk(lines)
	table.Render()
}
