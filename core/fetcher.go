package core

type IFetcher interface {
	IShutdown

	BeforeReq(task *Task)
	Fetch(task *Task)
	OnFetchSuccess(task *Task)
	OnFetchError(task *Task)

	Run() error
}
