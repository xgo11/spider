package core

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

const (
	DefaultPriority int = 5
	Utf8Encode          = "utf-8"
)

// Define task status
const (
	TaskStatusInit int = iota
	TaskStatusCrawled
	TaskStatusProcessed
	TaskStatusResulted
)

type TaskSchedule struct {
	Priority    int    `json:"priority"`
	ExecuteTime int64  `json:"execute_time" bson:"execute_time"`
	ITag        string `json:"i_tag" bson:"i_tag"`
	Force       bool   `json:"force" bson:"force"`
	AuthRecrawl bool   `json:"auth_recrawl" bson:"auth_recrawl"`
	Age         int64  `json:"age" bson:"age"`
}

type TaskFetcher struct {
	Method         string                 `json:"method" bson:"method"`
	Headers        map[string]string      `json:"headers" bson:"headers"`
	Cookies        map[string]string      `json:"cookies" bson:"cookies"`
	UseGzip        bool                   `json:"use_gzip" bson:"use_gzip"`
	Data           string                 `json:"data" bson:"data"`
	Proxy          string                 `json:"proxy" bson:"proxy"`
	Retries        int                    `json:"retries" bson:"retries"`
	MaxRedirects   int                    `json:"max_redirects" bson:"max_redirects"`
	ConnectTimeout int                    `json:"connect_timeout" bson:"connect_timeout"`
	Timeout        int                    `json:"timeout" bson:"timeout"`
	Save           map[string]interface{} `json:"save" bson:"save"`
}

type TaskProcessor struct {
	Callback       string `json:"callback" bson:"callback"`
	ProcessTimeout int    `json:"process_timeout" bson:"process_timeout"`
}

type Task struct {
	Url        string         `json:"url" bson:"url"`
	Project    string         `json:"project" bson:"project"`
	TaskId     string         `json:"task_id" bson:"task_id"`
	Catg       string         `json:"catg" bson:"catg"`
	SubCatg    string         `json:"sub_catg" bson:"sub_catg"`
	Status     int            `json:"status" bson:"status"`
	Schedule   *TaskSchedule  `json:"schedule" bson:"schedule"`
	Fetch      *TaskFetcher   `json:"fetch" bson:"fetch"`
	Process    *TaskProcessor `json:"process" bson:"process"`
	CreateTime int64          `json:"create_time" bson:"create_time"`
	UpdateTime int64          `json:"update_time" bson:"update_time"`
	LastCrawl  int64          `json:"last_crawl" bson:"last_crawl"`
}

type StatusMessage struct {
	TaskId string `json:"task_id"`
	Status int    `json:"status"`
}

type Schedule2FetchMessage struct {
	Task *Task `json:"task"`
}

type Fetch2ProcessMessage struct {
	Task     *Task     `json:"task"`
	Response *Response `json:"response"`
}

type Process2ResultMessage struct {
	Task   *Task   `json:"task"`
	Result *Result `json:"result"`
}

// fetch result process function
type ProcessFunc func(task *Task, response *Response, result *Result) error

type Response struct {
	StatusCode    int               `json:"status_code"`
	Url           string            `json:"url"`
	OrigUrl       string            `json:"orig_url"`
	Headers       http.Header       `json:"headers"`
	Cookies       map[string]string `json:"cookies"`
	Content       []byte            `json:"content"`
	ContentLength int               `json:"content_length"`
	TimeMS        int               `json:"time_ms"`
	ErrMessage    string            `json:"err_message"`
	Save          map[string]string `json:"save"`
	Encoding      string            `json:"encoding"`

	text string
	doc  *goquery.Document
}

func (r *Response) initEncoding() {

	if len(r.Encoding) > 0 {
		return
	}

	if r.Content != nil {
		if ctxLen := len(r.Content); ctxLen > 0 {
			contentType := r.Headers.Get("Content-Type")
			var data = r.Content
			if ctxLen >= 1024 {
				data = r.Content[:1024]
			}
			if _, name, ok := charset.DetermineEncoding(data, contentType); ok {
				r.Encoding = strings.Trim(name, " ")
			}
		}
	}

	if r.Encoding == "" {
		r.Encoding = Utf8Encode
	}

}

func (r *Response) IsOK() bool {
	if r.ErrMessage != "" {
		return false
	}
	if r.StatusCode != 200 {
		return false
	}
	return true
}

func (r *Response) GetEncoding() string {
	if r.Encoding == "" {
		r.initEncoding()
	}
	return r.Encoding
}

func (r *Response) GetText() string {
	if len(r.Content) == 0 {
		return ""
	}
	if r.text != "" {
		return r.text
	}

	if encoding, err := htmlindex.Get(r.GetEncoding()); err == nil {
		if name, _ := htmlindex.Name(encoding); name != Utf8Encode {
			if textBytes, err := encoding.NewDecoder().Bytes(r.Content); err == nil {
				if nil != textBytes && len(textBytes) > 0 {
					r.text = string(textBytes)
				}
			}
		}
	}

	if r.text == "" {
		r.text = string(r.Content)
	}
	return r.text
}

func (r *Response) Json(data interface{}) error {
	if err := json.Unmarshal(r.Content, data); err != nil {
		return err
	}
	return nil
}

func (r *Response) GetDocument() *goquery.Document {
	if r.doc == nil {
		if doc, err := goquery.NewDocumentFromReader(strings.NewReader(r.GetText())); err != nil {
			r.doc = doc
		}
	}
	return r.doc
}

type Result struct {
	ErrCode      int    `json:"err_code"`
	ErrMessage   string `json:"err_message"`
	Url          string `json:"url"`
	OrigUrl      string `json:"orig_url"`
	Html         string `json:"html"`
	NeedSnapshot bool   `json:"need_snapshot"`
	Parsed       []byte `json:"parsed"`
}
