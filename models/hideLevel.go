package models

type HideLevel string

const (
	Positive = "主动扫描"
	Passive  = "被动扫描"
)

var HideLevels []HideLevel = []HideLevel{
	Positive,
	Passive,
}

//DeceitWay 欺骗方式
type DeceitWay string

var (
	DeceitTarget  DeceitWay = "target"  //欺骗目标机
	DeceitGateway DeceitWay = "gateway" //欺骗网关
)

//PacketType ARP 包的类型
type PacketType string

var (
	Request PacketType = "request"
	Reply   PacketType = "reply"
)
