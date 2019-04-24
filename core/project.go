package core

import (
	"sync"
)

// fetch result process function
type ProcessCallback struct {
	Name     string
	Every    int64 // seconds
	Callback func(task *Task, response *Response) ([]*Task, *Result)
}

type IProject interface {
	GetName() string

	ListCallback() []ProcessCallback

	ListResultHook() []ResultWorkerHook
	ListProcessHook() []ProcessHook
	ListFetcherHook() []FetcherHook

	RegisterMe()
	ExecuteCallback(name string, task *Task, resp *Response) ([]*Task, *Result)
}

type IProjectBuilder interface {
	IProject
	AddCallback(callback ProcessCallback)
	AddProcessHook(hook ProcessHook)
	AddFetcherHook(hook FetcherHook)
	AddResultWorkerHook(hook ResultWorkerHook)
}

type IProjectManager interface {
	AddProject(project IProject)
	List() []IProject
	ListProjectNames() []string
	Get(name string) (IProject, bool)
}

type projectManager struct {
	sync.Mutex
	projects    []IProject
	projectsMap map[string]IProject
}

var (
	gProjectManager = &projectManager{}
)

func (pm *projectManager) AddProject(project IProject) {
	pm.Lock()
	defer pm.Unlock()

	var name = project.GetName()
	if pm.projectsMap == nil {
		pm.projectsMap = map[string]IProject{}
	}

	if _, ok := pm.projectsMap[name]; ok {
		return
	}

	pm.projectsMap[name] = project
	pm.projects = append(pm.projects, project)

}

func (pm *projectManager) List() []IProject {

	pm.Lock()
	defer pm.Unlock()
	return pm.projects

}

func (pm *projectManager) ListProjectNames() []string {
	pm.Lock()
	defer pm.Unlock()
	var names = make([]string, 0, len(pm.projects))
	for _, p := range pm.projects {
		names = append(names, p.GetName())
	}
	return names
}
func (pm *projectManager) Get(name string) (project IProject, ok bool) {
	pm.Lock()
	defer pm.Unlock()
	if pm.projectsMap != nil {
		project, ok = pm.projectsMap[name]
	}
	return
}

func GetProjectManager() IProjectManager {
	return gProjectManager
}
