package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/spf13/viper"
)

//Run 运行web服务器
func Run(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/index", homeHandler)
	mux.HandleFunc("/ws", wsHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("webspy.port")))
	if err != nil {
		log.Println("net.Listen failed,err:", err)
		return
	}
	go func() {
		err = http.Serve(listener, mux)
		if err != nil {
			fmt.Println("\r[*] 服务器关闭成功")
		}
	}()
	select {
	case <-ctx.Done():
	}
	listener.Close()
}
