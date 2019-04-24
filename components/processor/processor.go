package processor

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/xgo11/spider/common"
	"github.com/xgo11/spider/core"
)

type basicProcessor struct {
	sync.Mutex

	newTaskQ         core.IQueue
	fetcher2ProcessQ core.IQueue
	process2ResultQ  core.IQueue
	statusQ          core.IQueue
	pause            bool
	isRunning        bool
	wg               *sync.WaitGroup
	hooks            []core.ProcessHook
}

var (
	logger = common.GetLogger("processor")
)

func NewProcessor(newQ, f2pQ, p2rQ, sQ core.IQueue, hooks ...core.ProcessHook) core.IProcessor {
	p := &basicProcessor{}
	p.newTaskQ = newQ
	p.fetcher2ProcessQ = f2pQ
	p.process2ResultQ = p2rQ
	p.statusQ = sQ
	p.pause = false
	p.isRunning = false
	p.wg = &sync.WaitGroup{}
	p.hooks = make([]core.ProcessHook, 0, len(hooks))

	nameSet := make(map[string]bool)
	for _, h := range hooks {
		if !nameSet[h.Name] && (h.OnSendNewTask != nil || h.OnSendResult != nil) {
			p.hooks = append(p.hooks, h)
			nameSet[h.Name] = true
		}
	}

	return p
}

func (p *basicProcessor) Shutdown() {
	p.Lock()
	p.pause = true
	p.Unlock()
	p.wg.Wait()
	logger.Info("safe shutdown")
}

func (p *basicProcessor) HttpServe() http.HandlerFunc {
	engine := gin.Default()
	engine.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})

	engine.POST("/process", func(context *gin.Context) {

		var bodyBytes []byte
		var err error

		if bodyBytes, err = ioutil.ReadAll(context.Request.Body); err == nil {
			var body = core.Fetch2ProcessMessage{}
			if err = json.Unmarshal(bodyBytes, &body); err == nil {
				task, resp := body.Task, body.Response
				if task != nil && resp != nil {
					projectName := task.Project
					if project, exists := core.GetProjectManager().Get(projectName); exists {
						newTasks, result := project.ExecuteCallback(task.Process.Callback, task, resp)
						task.Status = core.TaskStatusProcessed
						output := struct {
							Task     *core.Task   `json:"task"`
							Result   *core.Result `json:"result"`
							NewTasks []*core.Task `json:"new_tasks"`
						}{
							Task:     task,
							Result:   result,
							NewTasks: newTasks,
						}
						context.JSON(http.StatusOK, &output)
						return
					} else {
						err = errors.New("project not exists")
					}
				}
			}
		}

		if err == nil {
			err = errors.New("invalid args")
		}

		context.String(http.StatusBadRequest, err.Error())
		return
	})
	return engine.ServeHTTP
}

func (p *basicProcessor) Run() {
	p.Lock()
	if p.isRunning {
		logger.Warnf("is running now!")
		p.Unlock()
		return
	}
	p.isRunning = true
	p.Unlock()

	var sleepIdle = 1000 * time.Millisecond

	logger.Info("start running ... ")

	for !p.pause {
		messages := p.fetcher2ProcessQ.Pop(1)
		if len(messages) < 1 && !p.pause {
			time.Sleep(sleepIdle)
		}
		for _, msg := range messages {
			body := core.Fetch2ProcessMessage{}
			if err := json.Unmarshal([]byte(msg), &body); err != nil || body.Task == nil || body.Response == nil {
				logger.WithError(err).Warn("invalid message")
				continue
			}
			go p.processOne(body.Task, body.Response)
		}
	}

	p.Lock()
	p.isRunning = false
	p.Unlock()
	logger.Info("stopped run")
}

func (p *basicProcessor) processOne(task *core.Task, resp *core.Response) {
	p.wg.Add(1)
	defer p.wg.Done()

	projectName := task.Project
	project, exists := core.GetProjectManager().Get(projectName)

	if !exists {
		logger.WithField("project", projectName).Warnf("project not exists")
		return
	}

	// TODO add callback timeout feature
	newTasks, result := project.ExecuteCallback(task.Process.Callback, task, resp)
	task.Status = core.TaskStatusProcessed

	if len(newTasks) > 0 { // send new tasks to scheduler
		for _, tsk := range newTasks {
			p.sendNewTask(tsk)
		}
	}

	if result != nil { // send result to result worker queue
		p.sendResult(task, result)
	}
}

func (p *basicProcessor) sendNewTask(task *core.Task) {
	if task != nil {
		task.Status = core.TaskStatusInit
		var err error
		var tskBytes []byte
		if tskBytes, err = json.Marshal(task); err == nil {
			if err = p.newTaskQ.Put(string(tskBytes)); err == nil {
				p.onSendNewTask(task)
			}
		}
		if err != nil {
			// TODO 兜底方案设计
			logger.WithError(err).WithField("task_info", string(tskBytes)).WithField("op", "sendNewTask").Error("fail")
		}
	}
}

func (p *basicProcessor) sendResult(task *core.Task, ret *core.Result) {
	msg := core.Process2ResultMessage{}
	msg.Task = task
	msg.Result = ret

	var msgBytes []byte
	var err error

	if msgBytes, err = json.Marshal(&msg); err == nil {
		p.onSendResult(task, ret)
	}
	if err != nil {
		// TODO 兜底方案
		logger.WithError(err).WithField("result", string(msgBytes)).WithField("op", "sendResult").Error("fail")
	}
}

func (p *basicProcessor) onSendNewTask(task *core.Task) {
	for _, h := range p.hooks {
		if h.OnSendNewTask != nil {
			h.OnSendNewTask(task)
		}
	}
}

func (p *basicProcessor) onSendResult(task *core.Task, ret *core.Result) {
	for _, h := range p.hooks {
		if h.OnSendResult != nil {
			h.OnSendResult(task, ret)
		}
	}
}
