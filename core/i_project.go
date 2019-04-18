package core

import (
	"github.com/xgo11/spider/model"
)

//定义项目接口
type IProject interface {
	IShutdown

	// 获取项目名称
	GetName() string

	//获取项目公用配置
	GetCrawlConfig() map[string]interface{}

	//获取项目回调处理库
	GetCallbacks() map[string]func(IHandler, *model.FetchResult) *model.ProcessResult //回调函数

	//项目启动函数，在生产环境下，一般实现为后台周期性队列处理过程
	OnStart()

	//注册项目到项目空间
	RegisterSelf()

	//生成一个handler实例
	GetHandler() IHandler

	//结果处理函数
	OnResult(*model.ProcessResult)

	GetResultWorker() IResultWorker

	GetTraceFields() []string
}
