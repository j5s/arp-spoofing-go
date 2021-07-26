package utils

import (
	"net"

	"github.com/malfunkt/iprange"
)

//GetIPList 解析多IP:string表示的ip 转化为 []net.IP 支持:10.0.0.1 192.168.1.0/24 192.168.1.0-255 192.168.1.* 四种形式
func GetIPList(ipStr string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ipStr)
	if err != nil {
		return nil, err
	}
	list := addressList.Expand()
	return list, err
}
