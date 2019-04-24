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

type Hook struct {
	Name    string
	Project string
}

type SchedulerHook struct {
	Hook
	OnTaskSelect func(task *Task)
	OnTaskNew    func(task *Task)
}

type IScheduler interface {
	IShutdown
	IRunnable
	IHTTPServer
}

type FetcherHook struct {
	Hook
	BeforeReq func(task *Task)
	AfterReq  func(task *Task, resp *Response)
}

type IFetcher interface {
	IShutdown
	IRunnable
	IHTTPServer
}

type ProcessHook struct {
	Hook
	OnSendNewTask func(*Task)
	OnSendResult  func(*Task, *Result)
}

type IProcessor interface {
	IShutdown
	IRunnable
	IHTTPServer
}

type ResultWorkerHook struct {
	Hook
	OnResult func(*Task, *Result)
}

type IResultWorker interface {
	IShutdown
	IRunnable
}
