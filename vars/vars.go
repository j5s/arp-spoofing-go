package vars

import (
	"context"

	"github.com/fatih/color"
)

var (
	//HostCancelMap 主机和通知退出函数之间的映射
	HostCancelMap map[string]context.CancelFunc = make(map[string]context.CancelFunc)
	//SniffCancelFunc 用于通知嗅探协程退出，如果为nil表示没有嗅探协程
	SniffCancelFunc context.CancelFunc = nil
	//MiddleAttackCancelMap 中间人攻击协程和通知退出之间的映射
	MiddleAttackCancelMap map[string]context.CancelFunc = make(map[string]context.CancelFunc)
	//Yellow 黄色输出
	Yellow = color.New(color.FgYellow).SprintFunc()
)
