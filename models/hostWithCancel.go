package models

import "context"

//HostWithCancel 可以退出的Host
type HostWithCancel struct {
	IP     string
	Cancel context.CancelFunc
}
