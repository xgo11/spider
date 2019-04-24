package core

import "net/http"

type IShutdown interface {
	Shutdown()
}

type IRunnable interface {
	Run()
}

type IHTTPServer interface {
	HttpServe() http.HandlerFunc
}

type IQueue interface {
	Name() string
	Put(message ...string) error
	Pop(count ...int) []string
	Size() int
	Limit() int
}

type SchedulerHook struct {
	Name         string
	OnTaskSelect func(task *Task)
	OnTaskNew    func(task *Task)
}

type IScheduler interface {
	IShutdown
	IRunnable
	IHTTPServer
}

type FetcherHook struct {
	Name      string
	BeforeReq func(task *Task)
	AfterReq  func(task *Task, resp *Response)
}

type IFetcher interface {
	IShutdown
	IRunnable
	IHTTPServer
}

type ProcessHook struct {
	Name          string
	OnSendNewTask func(*Task)
	OnSendResult  func(*Task, *Result)
}

type IProcessor interface {
	IShutdown
	IRunnable
	IHTTPServer
}
