package model

import (
	"strings"
	"time"
)

import (
	"github.com/xgo11/spider/utils"
)

const (
	KeyProjectName = "project"
	KeyAuthId      = "auth_id"
	KeyAuthTime    = "auth_time"
	KeyAuthChannel = "channel"
	KeyProxy       = "proxy" // 如果save中有此字段，使用指定的固定proxy抓取
	KeyCookies     = "cookies"

	KeyUserId      = "user_id"
	KeyUserProduct = "product"
	KeyTargetUid   = "target_uid"

	SchemeGoon  = "goon://"
	SchemeHttp  = "http://"
	SchemaHttps = "https://"

	CategoryPipeTimeout = "pipe_timeout" // 特殊类别，任务超时强制合并是使用
)

func NewCrawlTask() *CrawlTask {
	task := CrawlTask{CreateTime: time.Now().UnixNano(), Status: TaskStatusInitial}
	task.Schedule = &crawlTaskSchedule{Priority: TaskPriorityMiddle}
	task.Fetch = &crawlTaskFetch{Retries: 3, Method: "GET", MaxRedirects: 5, ConnectTimeout: 30, Timeout: 120}
	task.Process = &crawlTaskProcess{}

	task.CreateTime = time.Now().UnixNano() / 1e6

	return &task
}

func (t *CrawlTask) GetSave() map[string]interface{} {
	save := make(map[string]interface{})
	if nil != t.Fetch.Save {
		for k, v := range t.Fetch.Save {
			save[k] = v
		}
	}
	return save
}

func (t *CrawlTask) SetSave(save map[string]interface{}) {
	if len(save) > 0 {
		oldSave := t.GetSave()
		if nil == oldSave {
			oldSave = make(map[string]interface{})
		}
		for k, v := range save {
			oldSave[k] = v
		}
		t.Fetch.Save = oldSave
	}
}

func (t *CrawlTask) GetProxy() string {
	return t.Fetch.Proxy
}

func (t *CrawlTask) SetProxy(proxy string) {
	t.Fetch.Proxy = proxy
}

func (t *CrawlTask) GetMethod() string {
	return strings.ToUpper(t.Fetch.Method)
}

func (t *CrawlTask) GetPostData() string {
	return t.Fetch.Data
}

// 授权id
func (t *CrawlTask) GetAuthId() (result string) {
	return utils.FieldAsString(t.Fetch.Save, KeyAuthId)
}

// 授权时间
func (t *CrawlTask) GetAuthTime() (result string) {
	return utils.FieldAsString(t.Fetch.Save, KeyAuthTime)
}

func (t *CrawlTask) GetUserId() (result string) {
	return utils.FieldAsString(t.Fetch.Save, KeyUserId)
}

func (t *CrawlTask) GetUserProduct() (result string) {
	return utils.FieldAsString(t.Fetch.Save, KeyUserProduct)
}

func (t *CrawlTask) GetTargetUserId() (result string) {
	return utils.FieldAsString(t.Fetch.Save, KeyTargetUid)
}

func (t *CrawlTask) GetProjectName() string {
	return t.Project
}
