package project

import (
	"sync"
)
import (
	"github.com/xgo11/spider/core"
)

//ensure interface
var (
	_ core.IProject        = &standardProject{}
	_ core.IProjectBuilder = &standardProject{}
)

type standardProject struct {
	sync.Mutex

	name string

	cbMap map[string]core.ProcessCallback
	cbArr []core.ProcessCallback

	fHArr []core.FetcherHook
	pHArr []core.ProcessHook
	rHArr []core.ResultWorkerHook
}

func NewProjectBuilder(projectName string) core.IProjectBuilder {
	return &standardProject{name: projectName}
}

func (sp *standardProject) GetName() string {
	return sp.name
}

func (sp *standardProject) ListCallback() []core.ProcessCallback {
	return sp.cbArr
}

func (sp *standardProject) ListResultHook() []core.ResultWorkerHook {
	return sp.rHArr
}

func (sp *standardProject) ListProcessHook() []core.ProcessHook {
	return sp.pHArr
}

func (sp *standardProject) ListFetcherHook() []core.FetcherHook {
	return sp.fHArr
}

func (sp *standardProject) RegisterMe() {
	core.GetProjectManager().AddProject(sp)
}

func (sp *standardProject) ExecuteCallback(name string, task *core.Task, resp *core.Response) ([]*core.Task, *core.Result) {
	if sp.cbMap != nil {
		if cb, ok := sp.cbMap[name]; ok && cb.Callback != nil {
			return cb.Callback(task, resp)
		}
	}
	panic("no usable callback found")
}

func (sp *standardProject) AddCallback(callback core.ProcessCallback) {
	sp.Lock()
	defer sp.Unlock()

	if callback.Name == "" {
		panic("add empty name callback for project:" + sp.name)
	}

	if sp.cbMap == nil {
		sp.cbMap = map[string]core.ProcessCallback{}
	}

	if _, ok := sp.cbMap[callback.Name]; ok {
		panic("callback name conflict for project:" + sp.name)
	}

	sp.cbMap[callback.Name] = callback
	sp.cbArr = append(sp.cbArr, callback)

}

func (sp *standardProject) AddProcessHook(hook core.ProcessHook) {
	sp.Lock()
	defer sp.Unlock()

	nameSet := map[string]bool{}
	for _, h := range sp.pHArr {
		nameSet[h.Name] = true
	}

	if hook.Name == "" || nameSet[hook.Name] {
		panic("add invalid process hook for project:" + sp.name)
	}

	hook.Project = sp.name
	sp.pHArr = append(sp.pHArr, hook)

}

func (sp *standardProject) AddFetcherHook(hook core.FetcherHook) {
	sp.Lock()
	defer sp.Unlock()

	nameSet := map[string]bool{}
	for _, h := range sp.fHArr {
		nameSet[h.Name] = true
	}

	if hook.Name == "" || nameSet[hook.Name] {
		panic("add invalid fetcher hook for project:" + sp.name)
	}

	hook.Project = sp.name
	sp.fHArr = append(sp.fHArr, hook)

}

func (sp *standardProject) AddResultWorkerHook(hook core.ResultWorkerHook) {
	sp.Lock()
	defer sp.Unlock()

	nameSet := map[string]bool{}
	for _, h := range sp.rHArr {
		nameSet[h.Name] = true
	}

	if hook.Name == "" || nameSet[hook.Name] {
		panic("add invalid result worker hook for project:" + sp.name)
	}
	hook.Project = sp.name
	sp.rHArr = append(sp.rHArr, hook)

}
