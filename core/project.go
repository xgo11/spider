package core

import (
	"sync"
	"time"
)

// fetch result process function
type ProcessCallback struct {
	Name     string
	Every    time.Duration
	Callback func(task *Task, response *Response) ([]*Task, *Result)
}

type IProject interface {
	IShutdown
	GetName() string
	ListCallbacks() []string
	RegisterSelf()
	ExecuteCallback(name string, task *Task, resp *Response) ([]*Task, *Result)
}

type IProjectManager interface {
	IShutdown

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

func (pm *projectManager) Shutdown() {

	for _, p := range pm.projects {
		p.Shutdown()
	}

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
