package scheduler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xgo11/datetime"
	"github.com/xgo11/spider/common"
	"github.com/xgo11/spider/core"
)

type basicScheduler struct {
	sync.Mutex

	newTaskQ           core.IQueue
	scheduler2FetcherQ core.IQueue
	statusQ            core.IQueue
	hooks              []core.SchedulerHook
	pause              bool
	isRunning          bool
	wg                 *sync.WaitGroup
}

var (
	logger = common.GetLogger("scheduler")
)

func NewScheduler(newQ, s2fQ, sQ core.IQueue, hooks ...core.SchedulerHook) core.IScheduler {

	s := &basicScheduler{}
	s.newTaskQ = newQ
	s.scheduler2FetcherQ = s2fQ
	s.statusQ = sQ
	s.pause = false
	s.isRunning = false
	s.wg = new(sync.WaitGroup)

	nameSet := make(map[string]bool)
	s.hooks = make([]core.SchedulerHook, 0, len(hooks))
	for _, h := range hooks {
		if !nameSet[h.Name] && (h.OnTaskSelect != nil || h.OnTaskNew != nil) {
			s.hooks = append(s.hooks, h)
			nameSet[h.Name] = true
		}
	}
	return s
}

func (s *basicScheduler) Shutdown() {
	s.Lock()
	s.pause = true
	s.Unlock()

	s.wg.Wait()

	logger.Info("safe shutdown")

}

func (s *basicScheduler) Run() {

	s.Lock()
	if s.isRunning {
		logger.Warnf("is running, you should not run again")
		s.Unlock()
		return
	}
	s.isRunning = true
	s.Unlock()

	finC := make(chan struct{})

	var finished = 0
	projectList := core.GetProjectManager().List()
	var totalCount = 1 + len(projectList)

	go s.processTaskQueue(finC)
	for _, project := range projectList {
		go s.processProjectCron(project, finC)
	}

	for finished < totalCount {
		select {
		case <-finC:
			finished++
		}
	}

	s.Lock()
	s.isRunning = false
	s.Unlock()
	logger.Info("stopped run")
}

func (s *basicScheduler) processProjectCron(project core.IProject, fin chan struct{}) {
	defer func() {
		fin <- struct{}{}
	}()
	s.wg.Add(1)
	defer s.wg.Done()

	type st struct {
		lst   int64
		every int64
	}
	var crons = map[string]*st{}
	for _, cb := range project.ListCallbacks() {
		if cb.Every > 1 {
			crons[cb.Name] = &st{every: cb.Every, lst: 0}
		}
	}

	if len(crons) < 1 { // no cron jobs for this project
		return
	}

	projectName := project.GetName()
	logger.Infof("start crons for %v, count=%s", projectName, len(crons))

	for !s.pause {
		now := datetime.NowUnix()
		for name, cst := range crons {
			if now-cst.lst > cst.every {

				s.wg.Add(1)
				s.selectTask(&core.Task{
					Url:     core.SystemTaskSchema + "://" + name,
					Project: projectName,
					Status:  core.TaskStatusInit,
					Schedule: core.TaskSchedule{
						Priority: core.DefaultPriority,
					},
					Process: core.TaskProcessor{
						Callback: name,
					},
				})
				s.wg.Done()

				cst.lst = datetime.NowUnix()
			}
		}

		if !s.pause {
			time.Sleep(100 * time.Millisecond)
		}
	}

}
func (s *basicScheduler) processTaskQueue(fin chan struct{}) {
	defer func() {
		fin <- struct{}{}
	}()
	s.wg.Add(1)

	defer s.wg.Done()

	var oneBatchSize = 1000
	var sleepIdle = 1000 * time.Millisecond
	var sleepInterval = 100 * time.Millisecond

	for !s.pause {
		messages := s.newTaskQ.Pop(oneBatchSize)
		count := len(messages)
		if count > 0 {
			for _, msg := range messages {
				task := core.Task{}
				if err := json.Unmarshal([]byte(msg), &task); err != nil {
					logger.WithError(err).WithFields(logrus.Fields{
						"task": msg,
						"op":   "abandon_new",
					}).Warn("invalid task message")
					continue
				}
				go s.receiveNewTask(&task)
			}

		}
		if !s.pause {
			if count < oneBatchSize {
				time.Sleep(sleepIdle)
			} else {
				time.Sleep(sleepInterval)
			}
		}
	}
}

func (s *basicScheduler) receiveNewTask(task *core.Task) {
	s.wg.Add(1)
	defer s.wg.Done()
	s.onReceiveNew(task)
	//TODO 调度算法优化
	s.selectTask(task)

}

func (s *basicScheduler) selectTask(task *core.Task) {
	task.Status = core.TaskStatusScheduled
	msg := core.Schedule2FetchMessage{Task: task}

	var msgBytes []byte
	var err error

	if msgBytes, err = json.Marshal(&msg); err == nil {
		if err = s.scheduler2FetcherQ.Put(string(msgBytes)); err == nil {
			s.onSelect(task)
		}
	}

	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"op":   "selectTask",
			"task": string(msgBytes),
		}).Error("fail")
	}
}
func (s *basicScheduler) onReceiveNew(task *core.Task) {
	for _, h := range s.hooks {
		if h.OnTaskNew != nil {
			h.OnTaskNew(task)
		}
	}
}

func (s *basicScheduler) onSelect(task *core.Task) {

	for _, h := range s.hooks {
		if h.OnTaskSelect != nil {
			h.OnTaskSelect(task)
		}
	}
}

func (s *basicScheduler) HttpServe() http.HandlerFunc {
	eng := gin.Default()
	eng.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})
	return eng.ServeHTTP
}
