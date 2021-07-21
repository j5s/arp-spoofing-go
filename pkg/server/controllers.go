package server

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

//IndexData 页面需要的数据
type IndexData struct {
	Data string
	Host string
}

//homeHandler  首页
func homeHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析模板生成模板对象
	tmpl, err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	//2.填充数据，渲染模板
	indexData := IndexData{
		Data: "",
		Host: r.Host,
	}
	tmpl.Execute(w, indexData)
}

//wsHandler 处理websocket
func wsHandler(w http.ResponseWriter, r *http.Request) {
	//1.批准客户端升级协议为ws,和客户端建立连接
	upgrader := websocket.Upgrader{
		ReadBufferSize:  10240,
		WriteBufferSize: 10240,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader.Upgrade failed,err:", err)
		return
	}
	//2.向客户端写
	go writeTo(conn)
	//3.从客户端读
	readFrom(conn)
}
