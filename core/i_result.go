package core

import (
	"github.com/xgo11/spider/model"
)

//定义项目接口
type IResultWorker interface {
	IShutdown

	// 获取结果db实例
	ResultDB() interface{}

	// 结果合并入口
	OnResult(*model.ProcessResult) error

	// 周期性强制合并任务
	CronOnResult()

	// 结果通知
	RunNoticeChain(map[string]interface{})

	// add Notification
	AddNotification(INotification) bool
}

type INotification interface {
	GetName() string
	Notice(message map[string]interface{}) (string, bool)
}
