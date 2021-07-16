package settings

import (
	"ARPSpoofing/models"
	"ARPSpoofing/utils"
	"ARPSpoofing/vars"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

//ARPScanOptions ARP扫描配置
var (
	Options models.Options
)

//Init 初始化配置
func Init() (err error) {
	//指定配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/arp/")
	viper.AddConfigPath("./conf/")

	//读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed,err:%v\n", err)
		return err
	}
	options := models.NewOptions("ARP Scan Options")

	//获取所有网卡
	vars.IfaceNames, err = utils.GetAllIfaceNames(true)
	if err != nil {
		log.Println("Utils.GetAllIfaceNames(true) failed,err:", err)
		return err
	}
	//初始化扫描配置项
	options.Add("ifname", vars.IfaceNames[0], true, "监听哪个网卡")
	options.Add("range", "192.168.1.0/24", true, "扫描范围")
	options.Add("type", "all", true, "扫描方式")
	return nil
}
