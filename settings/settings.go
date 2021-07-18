package settings

import (
	"ARPSpoofing/models"
	"ARPSpoofing/utils"
	"ARPSpoofing/vars"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/spf13/viper"
)

//ARPScanOptions ARP扫描配置
var (
	Options *models.Options
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
	Options = models.NewOptions("ARP Scan Options")

	//获取默认网卡
	vars.IfaceNames, err = utils.GetAllIfaceNames(true)
	if err != nil {
		log.Println("Utils.GetAllIfaceNames(true) failed,err:", err)
		return err
	}
	iface, err := net.InterfaceByName(vars.IfaceNames[0])
	if err != nil {
		log.Println(err)
		return err
	}
	//获取默认扫描范围
	myIP, err := utils.GetIPv4ByIface(iface)
	if err != nil {
		log.Println(err)
		return err
	}
	scanRange := fmt.Sprintf("%s/24", myIP.String())
	gateway := fmt.Sprintf("%s.1", strings.Join(strings.Split(myIP.String(), ".")[:3], "."))
	//初始化扫描配置项
	Options.Add("ifname", vars.IfaceNames[0], true, "监听哪个网卡")
	Options.Add("range", scanRange, true, "扫描范围")
	Options.Add("method", "all", true, "扫描方式:all,arp,udp")
	Options.Add("gateway", gateway, true, "局域网的网关")
	return nil
}
