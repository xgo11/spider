package core

import "net/http"

type ProcessHook struct {
	Name          string
	OnSendNewTask func(*Task)
	OnSendResult  func(*Task, *Result)
}

type IProcessor interface {
	IShutdown
	Run()
	HttpServe() http.HandlerFunc
}
