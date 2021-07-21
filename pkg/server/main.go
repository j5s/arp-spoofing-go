package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

//Run 运行web服务器
func Run() {
	http.HandleFunc("/index", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("webspy.port")), nil)
		if err != nil {
			log.Println("http.ListenAndServe failed,err:", err)
			return
		}
	}()
	log.Printf("启动服务器:http://127.0.0.1:%d/index", viper.GetInt("webspy.port"))
}
