package spider

import (
	"github.com/xgo11/spider/core"
	"github.com/xgo11/spider/curl"
)

func UrlTask(url string, kwArgs map[string]interface{}) *core.Task {
	tsk := core.NewTask(url)
	tsk.Update(kwArgs)
	return tsk
}

func CurlTask(curlCommand string) (task *core.Task, err error) {
	if url, kwArgs, err := curl.AnalyzeCurl(curlCommand); err != nil {
		return nil, err
	} else {
		task = UrlTask(url, kwArgs)
		return
	}
}

func BuildResult(resp *core.Response) *core.Result {
	result := &core.Result{
		ErrCode:    resp.StatusCode,
		ErrMessage: resp.ErrMessage,
		Url:        resp.Url,
		OrigUrl:    resp.OrigUrl,
		Html:       resp.GetText(),
	}
	return result
}
