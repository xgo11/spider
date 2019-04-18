package core

import (
	"github.com/xgo11/spider/model"
)

//定义任务处理基础结构，主要职责如下：
// - 提供提交任务基础函数，发送任务到处理队列
// - 提供任务结果解析函数挂载点，便于网络请求完成后进行回调
//
type IHandler interface {
	IShutdown

	//
	//基础信息
	GetProjectName() string
	GetCrawlConfig() map[string]interface{}
	SetCrawlConfigKv(k string, v interface{}) IHandler
	SetCrawlConfig(m map[string]interface{}) IHandler

	//动态绑定任务
	SetCurrentTask(task *model.CrawlTask) bool
	CurrentTask() *model.CrawlTask

	//挂载回调处理函数
	RegisterCallback(string, func(IHandler, *model.FetchResult) *model.ProcessResult) bool
	RunCallback(string, *model.FetchResult) (*model.ProcessResult, error)

	//创建任务
	Crawl(string, map[string]interface{}) ([]*model.CrawlTask, error)

	GetTaskId(task *model.CrawlTask) string
}
