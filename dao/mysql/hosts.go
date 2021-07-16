package mysql

import (
	"ARPSpoofing/models"
	"log"
)

//AddHost 添加主机
func AddHost(host *models.Host) error {
	sql := "insert into hosts(ip,mac,mac_info) values(?,?,?)"
	_, err := db.Exec(sql, host.IP, host.MAC, host.MACInfo)
	if err != nil {
		log.Println("db.Exec(sql, host.IP, host.MAC, host.MACInfo) failed,err:", err)
		return err
	}
	return nil
}
