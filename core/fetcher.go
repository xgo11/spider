package core

import "net/http"

type FetcherHook struct {
	Name      string
	BeforeReq func(task *Task)
	AfterReq  func(task *Task, resp *Response)
}

type IFetcher interface {
	IShutdown
	Run() error

	HttpServe() http.HandlerFunc
}
