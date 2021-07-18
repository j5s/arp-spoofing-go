package vars

import "context"

//ARPScanOptions ARP扫描配置
var (
	IfaceNames []string
)

//HostCancelMap 主机和通知退出函数之间的映射
var HostCancelMap map[string]context.CancelFunc = make(map[string]context.CancelFunc)

//SniffCancelFunc 用于通知嗅探协程退出，如果为nil表示没有嗅探协程
var SniffCancelFunc context.CancelFunc = nil
