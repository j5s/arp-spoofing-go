package logic

import (
	"ARPSpoofing/dao/redis"
	"fmt"
	"log"
)

//ShowLoot 展示所有战利品
func ShowLoot() error {
	loots, err := redis.NewLootsSaver().GetAll()
	if err != nil {
		log.Println("redis get all failed,err:", err)
	}
	if len(loots) == 0 {
		return nil
	}
	for _, loot := range loots {
		fmt.Printf("%s:%s->%s:%s [%s]\n", loot.SrcIP, loot.SrcPort, loot.DstIP, loot.DstPort, loot.Keyword)
		fmt.Println(loot.Payload)
	}
	return nil
}

//ClearLoot 清除所有战利品
func ClearLoot() error {
	if err := redis.NewLootsSaver().ClearAll(); err != nil {
		return err
	}
	return nil
}
