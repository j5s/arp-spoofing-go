package logic

import (
	"ARPSpoofing/dao/redis"
	"ARPSpoofing/pkg/table"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/fatih/structs"
)

//ShowHosts 展示所有主机
func ShowHosts() error {
	hosts, err := redis.NewHosts().GetAll()
	if err != nil {
		log.Println(err)
		return err
	}
	if len(hosts) == 0 {
		fmt.Println(color.YellowString("[*] 没有主机,请先 scan 扫描局域网"))
		return nil
	}
	headers := structs.Names(hosts[0])
	data := make([][]string, 0, len(hosts))
	for index, host := range hosts {
		data = append(data, []string{
			fmt.Sprintf("%v", index),
			host.IP,
			host.MAC,
			host.MACInfo,
		})
	}
	table.Show(headers, data)
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
