package server

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

//writeTo 向Conn中写
func writeTo(conn *websocket.Conn) {
	for range time.Tick(time.Second * 2) {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("happy hacking")); err != nil {
			log.Println("conn write failed,err:", err)
			return
		}
	}
}

//readFrom 从Conn中读
func readFrom(conn *websocket.Conn) {
	return
}
