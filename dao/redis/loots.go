package redis

import (
	"ARPSpoofing/models"
	"encoding/json"
	"log"
)

//LootsSaver 战利品 敏感信息
type LootsSaver struct {
	key string
}

//NewLootsSaver 新建一个存储器
func NewLootsSaver() *LootsSaver {
	return &LootsSaver{
		key: "loots",
	}
}

//Add 添加战利品
func (l *LootsSaver) Add(loot *models.Loot) error {
	lootStr, err := json.Marshal(loot)
	if err != nil {
		log.Println("json.Marshal failed,err:", err)
		return err
	}
	_, err = rdb.SAdd(l.key, lootStr).Result()
	if err != nil {
		log.Println("rdb.SAdd failed,err:", err)
		return err
	}
	return nil
}

//GetAll 获取所有战利品
func (l *LootsSaver) GetAll(loot *models.Loot) ([]models.Loot, error) {
	items, err := rdb.SMembers(l.key).Result()
	if err != nil {
		return nil, err
	}
	loots := make([]models.Loot, 0, len(items))
	var temp models.Loot
	for index := range items {
		err = json.Unmarshal([]byte(items[index]), &temp)
		if err != nil {
			return nil, err
		}
		loots = append(loots, temp)
	}
	return loots, nil
}
