package socket

import (
	"net"
)



type HideLevel string

const (
	Positive = "主动扫描"
	Passive  = "被动扫描"
)

var HideLevels []HideLevel = []HideLevel{
	Positive,
	Passive,
}
