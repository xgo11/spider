package model

import (
	"encoding/json"
	"net/http"
	"strings"
)

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"
)

//定义网络请求结果
type FetchResult struct {
	StatusCode    int               `json:"status_code"`    // http状态码
	Url           string            `json:"url"`            // 最终url
	OrigUrl       string            `json:"orig_url"`       // 原始url
	Headers       http.Header       `json:"headers"`        // 最终header
	Cookies       map[string]string `json:"cookies"`        // 最终cookie，经过合并
	Content       []byte            `json:"content"`        // 最终返回数据内容
	ContentLength int64             `json:"content_length"` // 返回数据长度
	Error         string            `json:"error"`          // 错误信息
	CostTime      int64             `json:"cost_time"`      // 精确到毫秒
	Task          *CrawlTask        `json:"task"`           // 原始任务信息

	//解析信息
	encoding string            // 内容编码格式
	text     string            // 文本内容
	doc      *goquery.Document // dom
}

func (r *FetchResult) IsOk() bool {
	if len(r.Error) > 0 {
		return false
	}
	return true
}

func (r *FetchResult) TryGetJson() interface{} {
	var result interface{}
	if err := json.Unmarshal([]byte(r.GetText()), result); err != nil {
		return nil
	}
	return result
}

func (r *FetchResult) TryGetGoQuery() (doc *goquery.Document) {
	if r.doc == nil {
		r.doc, _ = goquery.NewDocumentFromReader(strings.NewReader(r.GetText()))
	}
	doc = r.doc
	return
}

// 获取UTF-8编码格式的文本
func (r *FetchResult) GetText() string {
	if r.Content == nil || len(r.Content) == 0 {
		return ""
	}

	if len(r.text) > 0 {
		return r.text
	}

	charEncoding := r.GetEncoding()
	if e, err := htmlindex.Get(charEncoding); err == nil {
		if name, _ := htmlindex.Name(e); name != "utf-8" {
			if textBytes, err := e.NewDecoder().Bytes(r.Content); err == nil {
				if nil != textBytes && len(textBytes) > 0 {
					r.text = string(textBytes)
					return r.text
				}
			}
		} else {
			r.text = string(r.Content)
			return r.text
		}
	}
	return string(r.Content)
}

//获取源编码方式
func (r *FetchResult) GetEncoding() string {

	if len(r.encoding) > 0 {
		return r.encoding
	}

	if r.Content != nil {
		if ctxLen := len(r.Content); ctxLen > 0 {
			contentType := r.Headers.Get("Content-Type")
			var data = r.Content
			if ctxLen >= 1024 {
				data = r.Content[:1024]
			}
			if _, name, ok := charset.DetermineEncoding(data, contentType); ok {
				r.encoding = strings.Trim(name, " ")
				return r.encoding
			}
		}
	}
	return "utf-8"
}
