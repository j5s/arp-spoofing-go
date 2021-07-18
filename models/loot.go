package models

//Loot 含有敏感信息的报文（战利品）
type Loot struct {
	SrcIP   string
	DstIP   string
	SrcPort string
	DstPort string
	Payload string
	Keyword string
}
