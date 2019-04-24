package fetcher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xgo11/spider/common"
	"github.com/xgo11/spider/core"
)

type httpFetcher struct {
	sync.Mutex

	schedule2FetcherQ  core.IQueue
	fetcher2ProcessorQ core.IQueue
	statusQ            core.IQueue
	hooks              []core.FetcherHook
	wg                 *sync.WaitGroup
	pause              bool
	isRunning          bool
}

const (
	sleepIdle = 3000 * time.Millisecond
	loopSize  = 20
)

var (
	logger = common.GetLogger("fetcher")
)

func (hf *httpFetcher) Shutdown() {
	hf.Lock()
	hf.pause = true
	hf.Unlock()
	hf.wg.Wait()
	logger.Infof("safe stopped")
}

func (hf *httpFetcher) Run() {
	hf.Lock()
	if hf.isRunning {
		hf.Unlock()
		logger.Warnf("fetcher is running, you should not run again")
		return
	}
	hf.isRunning = true
	hf.Unlock()

	logger.Infof("starting ...")

	var oneLoopSize = loopSize

	for !hf.pause {
		messages := hf.schedule2FetcherQ.Pop(oneLoopSize)
		count := len(messages)

		if count > 0 {
			for _, msg := range messages {
				body := core.Schedule2FetchMessage{}
				if err := json.Unmarshal([]byte(msg), &body); err != nil || body.Task == nil {
					logger.WithError(err).WithField("op", "abandon").Error(msg[0:300])
					continue
				}
				go hf.runOneTask(body.Task)
			}
		}

		if count < oneLoopSize {
			time.Sleep(sleepIdle)
		}
	}

	hf.Lock()
	hf.isRunning = false
	logger.Infof("stopped run loop")
	hf.Unlock()
	return
}

func (hf *httpFetcher) HttpServe() http.HandlerFunc {
	//gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
	engine.POST("/fetch", func(ctx *gin.Context) {
		var reqBytes []byte
		var err error

		if reqBytes, err = ioutil.ReadAll(ctx.Request.Body); err == nil {
			task := core.Task{}
			if err = json.Unmarshal(reqBytes, &task); err == nil {
				resp := hf.fetch(&task)
				resp.GetEncoding()
				out := core.Fetch2ProcessMessage{Task: &task, Response: resp}
				ctx.JSON(http.StatusOK, out)
			}
		}

		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		}
	})
	return engine.ServeHTTP
}

func (hf *httpFetcher) runOneTask(task *core.Task) {
	hf.wg.Add(1)
	defer hf.wg.Done()

	if resp := hf.fetch(task); resp != nil {
		hf.onSendMessage(task, resp)
	}

}

func (hf *httpFetcher) fetch(task *core.Task) (resp *core.Response) {

	var uri *url.URL
	var err error

	if uri, err = url.Parse(task.Url); err != nil || uri == nil || uri.Scheme == "" {
		logger.WithError(err).WithFields(logrus.Fields{
			"op":      "abandon",
			"taskid":  task.TaskId,
			"project": task.Project,
			"url":     task.Url,
		}).Error("invalid url")
		return
	}

	if uri.Scheme == core.SystemTaskSchema { //system schedule task
		resp = &core.Response{StatusCode: 200, Url: task.Url}
		return
	}

	// http requests
	hf.beforeReq(task)
	resp = (&httpClient{}).Do(task)

	if resp.ErrMessage == "" {
		hf.onFetchSuccess(task, resp)
	} else {
		hf.onFetchError(task, resp)
	}

	hf.afterReq(task, resp)
	return

}

func (hf *httpFetcher) registerHooks(hooks ...core.FetcherHook) {
	nameSet := map[string]bool{}
	for _, h := range hf.hooks {
		nameSet[h.Name] = true
	}

	for _, h := range hooks {
		if nameSet[h.Name] || (h.BeforeReq == nil && h.AfterReq == nil) {
			continue
		}
		hf.hooks = append(hf.hooks, h)
		nameSet[h.Name] = true
	}
}

func (hf *httpFetcher) onFetchSuccess(task *core.Task, resp *core.Response) {
	logger.WithFields(logrus.Fields{
		"taskid":      task.TaskId,
		"url":         task.Url,
		"cost":        resp.TimeMS,
		"status_code": resp.StatusCode,
	}).Info("ok")
}

func (hf *httpFetcher) onFetchError(task *core.Task, resp *core.Response) {
	logger.WithFields(logrus.Fields{
		"taskid":      task.TaskId,
		"url":         task.Url,
		"cost":        resp.TimeMS,
		"status_code": resp.StatusCode,
		"error":       resp.ErrMessage,
		"proxy":       task.Fetch.Proxy,
	}).Error("fail")
}

func (hf *httpFetcher) onSendMessage(task *core.Task, resp *core.Response) {
	var bytes []byte
	var err error
	if bytes, err = json.Marshal(core.Fetch2ProcessMessage{Task: task, Response: resp}); err == nil {
		err = hf.fetcher2ProcessorQ.Put(string(bytes))
	}
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"op":      "sendMessage",
			"taskid":  task.TaskId,
			"project": task.Project,
		}).Error()
	}
}

func (hf *httpFetcher) beforeReq(task *core.Task) {
	for _, h := range hf.hooks {
		if h.BeforeReq != nil {
			h.BeforeReq(task)
		}
	}

	if proj, ok := core.GetProjectManager().Get(task.Project); ok {
		for _, h := range proj.ListFetcherHook() {
			if h.BeforeReq != nil {
				h.BeforeReq(task)
			}
		}
	}

}

func (hf *httpFetcher) afterReq(task *core.Task, resp *core.Response) {
	for _, h := range hf.hooks {
		if h.AfterReq != nil {
			h.AfterReq(task, resp)
		}
	}

	if proj, ok := core.GetProjectManager().Get(task.Project); ok {
		for _, h := range proj.ListFetcherHook() {
			if h.AfterReq != nil {
				h.AfterReq(task, resp)
			}
		}
	}
}
