package result_worker

import (
	"encoding/json"
	"sync"
	"time"
)

import (
	"github.com/sirupsen/logrus"
	"github.com/xgo11/spider/common"
	"github.com/xgo11/spider/core"
)

var (
	logger = common.GetLogger("result_worker")
)

type basicResultWorker struct {
	sync.Mutex

	process2ResultQ core.IQueue
	statusQ         core.IQueue
	wg              *sync.WaitGroup
	hooks           []core.ResultWorkerHook
	hooksCount      int
	pause           bool
	isRunning       bool
}

func NewResultWorker(p2rQ, sQ core.IQueue, hooks ...core.ResultWorkerHook) core.IResultWorker {
	r := &basicResultWorker{}
	r.process2ResultQ = p2rQ
	r.statusQ = sQ
	r.wg = new(sync.WaitGroup)
	r.pause = false
	r.isRunning = false

	nameSet := make(map[string]bool)
	r.hooks = make([]core.ResultWorkerHook, 0, len(hooks))

	for _, h := range hooks {
		if !nameSet[h.Name] && h.OnResult != nil {
			r.hooks = append(r.hooks, h)
			nameSet[h.Name] = true
		}
	}
	r.hooksCount = len(r.hooks)
	return r
}

func (r *basicResultWorker) Shutdown() {
	r.Lock()
	r.pause = true
	r.Unlock()
	r.wg.Wait()

	logger.Info("safe shutdown")
}

func (r *basicResultWorker) Run() {
	r.Lock()
	if r.isRunning {
		logger.Warnf("is running, you should not run again")
		r.Unlock()
		return
	}
	r.isRunning = true
	r.Unlock()

	logger.Info("running ... ")

	for !r.pause {
		if messages := r.process2ResultQ.Pop(1); len(messages) < 1 {
			time.Sleep(2 * time.Second)
		} else {
			body := core.Process2ResultMessage{}
			if err := json.Unmarshal([]byte(messages[0]), &body); err == nil && body.Task != nil && body.Result != nil {
				go r.onResult(body.Task, body.Result)
			} else {
				// TODO 降级方案
				logger.WithError(err).WithField("body", messages[0]).WithField("op", "onResult").Error("fail")
			}
		}
	}

	r.Lock()
	r.isRunning = false
	r.Unlock()
	logger.Info("stopped run")
}

func (r *basicResultWorker) onResult(task *core.Task, ret *core.Result) {

	task.Status = core.TaskStatusResulted

	if r.hooksCount < 1 { // when there is no hooks
		logger.WithFields(logrus.Fields{
			"taskid": task.TaskId,
			"url":    task.Url,
			"code":   ret.ErrCode,
			"err":    ret.ErrMessage,
			"ret":    string(ret.Parsed),
		}).Info("show")
	}

	for _, h := range r.hooks {
		if h.OnResult != nil {
			h.OnResult(task, ret)
		}
	}

}
