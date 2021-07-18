package logic

import (
	"ARPSpoofing/dao/redis"
	"fmt"
	"log"
	"os"

	"github.com/fatih/structs"
	"github.com/olekukonko/tablewriter"
)

//ShowHosts 展示所有主机
func ShowHosts() error {
	hosts, err := redis.NewHosts().GetAll()
	if err != nil {
		log.Println(err)
		return err
	}
	if len(hosts) == 0 {
		return nil
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoMergeCells(true)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetHeader(structs.Names(hosts[0]))

	for index, host := range hosts {
		table.Append([]string{
			fmt.Sprintf("%v", index),
			host.IP,
			host.MAC,
			host.MACInfo,
		})
	}
	table.Render() // Send output
	return nil
}

//ClearHosts  清除所有主机
func ClearHosts() error {
	err := redis.NewHosts().Clear()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
