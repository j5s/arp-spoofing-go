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

//IsHostExist 判断主机是否存在
func IsHostExist(host *models.Host) (bool, error) {
	sql := "select count(*) from hosts where ip=?"
	var count int
	err := db.Get(&count, sql, host.IP)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//GetAllHosts 获取所有主机信息
func GetAllHosts() ([]models.Host, error) {
	sql := `select * from hosts`
	var hosts []models.Host
	err := db.Select(&hosts, sql)
	if err != nil {
		log.Println("db.Select hosts failed,err:", err)
		return nil, err
	}
	return hosts, nil
}

//ClearHosts 清空所有主机
func ClearHosts() error {
	sql := `delete from hosts`
	_, err := db.Exec(sql)
	if err != nil {
		log.Println("clear hosts failed,err:", err)
		return err
	}
	return nil
}
