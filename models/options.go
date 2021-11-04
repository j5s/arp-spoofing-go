package models

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/olekukonko/tablewriter"
)

//OptionItem 配置项
type OptionItem struct {
	Name        string
	Value       string
	Required    bool
	Description string
}

//Options 配置
type Options struct {
	Name    string
	Content []OptionItem
}

//NewOptions 创建配置对象
func NewOptions(name string) *Options {
	return &Options{
		Name: name,
	}
}

//Add 添加配置项
func (o *Options) Add(name string, value string, required bool, description string) {
	o.Content = append(o.Content, OptionItem{
		Name:        name,
		Value:       value,
		Required:    required,
		Description: description,
	})
}

//Set 设置配置项
func (o *Options) Set(name string, value string) {
	for i := range o.Content {
		if o.Content[i].Name == name {
			o.Content[i].Value = strings.TrimSpace(value)
		}
	}
}

//Get 获取配置项的值
func (o *Options) Get(name string) (string, error) {
	for i := range o.Content {
		if o.Content[i].Name == name {
			return o.Content[i].Value, nil
		}
	}
	return "", fmt.Errorf("key %s not exist", name)
}

//Show 展示配置项
func (o *Options) Show() {
	lines := make([][]string, 0)
	for _, option := range o.Content {
		line := make([]string, 0, 3)
		for _, value := range structs.Values(option) {
			line = append(line, fmt.Sprintf("%v", value))
		}
		lines = append(lines, line)
	}
	var optionItem OptionItem
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(structs.Names(optionItem))
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.AppendBulk(lines)
	table.SetCaption(true, o.Name)
	table.Render()
}
