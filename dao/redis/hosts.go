package redis

import (
	"ARPSpoofing/models"
	"encoding/json"
	"log"

	"github.com/sirupsen/logrus"
)

type Hosts struct {
	key string
}

//NewHosts 扫描到的所有主机
func NewHosts() *Hosts {
	return &Hosts{
		key: "hosts",
	}
}

//Add 添加主机
func (h *Hosts) Add(host *models.Host) error {
	hoststr, err := json.Marshal(host)
	if err != nil {
		log.Println("json.Marshal(host) faild,err:", err)
		return err
	}
	_, err = rdb.HSet(h.key, host.IP, hoststr).Result()
	if err != nil {
		log.Println("rdb.HSet failed,err:", err)
		return err
	}
	return nil
}

//Get 获取一个主机的详细信息
func (h *Hosts) Get(ip string) (models.Host, error) {
	var host models.Host
	ret, err := rdb.HGet(h.key, ip).Result()
	if err != nil {
		logrus.Errorf("rdb.HGet(%v, %v).Result() faild,err:%v\n", h.key, ip, err)
		return host, err
	}
	err = json.Unmarshal([]byte(ret), &host)
	if err != nil {
		logrus.Errorf("json.Unmarshal([]byte(ret), &host) failed,err:%v\n", err)
		return host, err
	}
	return host, nil
}

func (h *Hosts) GetAllIP() ([]string, error) {
	return rdb.HKeys(h.key).Result()
}

func (h *Hosts) GetAll() ([]models.Host, error) {
	hosts := make([]models.Host, 0)
	ret, err := rdb.HGetAll(h.key).Result()
	if err != nil {
		return nil, err
	}
	var temp models.Host
	for _, value := range ret {
		err := json.Unmarshal([]byte(value), &temp)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, temp)
	}
	return hosts, nil
}

func (h *Hosts) Clear() error {
	_, err := rdb.Del(h.key).Result()
	if err != nil {
		return err
	}
	return nil
}
