package model

const (
	// 任务优先级，优先级值大的优先调度
	TaskPriorityMinimum = 1
	TaskPriorityMaximum = 9
	TaskPriorityMiddle  = 5

	// 任务状态
	TaskStatusInitial  = 0 // 初始
	TaskStatusFetched  = 1 // 网络请求完成
	TaskStatusParsed   = 2 // 解析完成
	TaskStatusResulted = 3 // 结果入库完成

)

var (
	TaskFieldsForSchedule = [...]string{
		"priority",
		"exetime",
		"age",
		"itag",
		"force_update",
		"auto_recrawl",
	}

	TaskFieldsForFetch = [...]string{
		"method",
		"headers",
		"user_agent",
		"cookies",
		"use_gzip",
		"etag",
		"last_modifed",
		"last_modified",
		"data",
		"proxy",
		"retries",
		"max_redirects",
		"connect_timeout",
		"timeout",
		"save",
		"js_run_at",
		"js_script",
		"js_viewport_width",
		"js_viewport_height",
		"load_images",
		"fetch_type",
		"validate_cert",
		"robots_txt",
	}

	TaskFieldsForProcess = [...]string{
		"callback",
		"process_time_limit",
	}
)

//抓取任务数据结构
type CrawlTask struct {
	Url     string `json:"url" bson:"url"`           // 任务url
	Project string `json:"project" bson:"project"`   // 任务所属的项目名称
	TaskId  string `json:"task_id" bson:"task_id"`   // 任务唯一id
	Catg    string `json:"catg" bson:"catg"`         // 任务分类
	SubCatg string `json:"sub_catg" bson:"sub_catg"` // 任务二级分类
	Status  int    `json:"status" bson:"status"`     // 任务当前状态

	CreateTime    int64 `json:"create_time" bson:"create_time"`         //精确到毫秒
	UpdateTime    int64 `json:"update_time" bson:"update_time"`         //精确到毫秒
	LastCrawlTime int64 `json:"last_crawl_time" bson:"last_crawl_time"` //精确到毫秒

	Schedule *crawlTaskSchedule `json:"schedule" bson:"schedule"` // 调度相关参数
	Fetch    *crawlTaskFetch    `json:"fetch" bson:"fetch"`       // 网络请求相关参数
	Process  *crawlTaskProcess  `json:"process" bson:"process"`   // 处理相关参数
}

// 任务调度相关参数
//"priority",
//"exetime",
//"age",
//"itag",
//"force_update",
//"auto_recrawl"
type crawlTaskSchedule struct {
	Priority    int    `json:"priority" bson:"priority"`         // 任务优先级
	ExecuteTime int64  `json:"exetime" bson:"exetime"`           // 任务执行时间,Unix时间戳
	ITag        string `json:"itag" json:"itag"`                 // 任务标签
	ForceUpdate bool   `json:"force_update" bson:"force_update"` // 是否强制更新

	// 自动重抓
	AutoRecrawl bool  `json:"auto_recrawl" bson:"auto_recrawl"` // 是否启用自动周期性重抓
	Age         int64 `json:"age" bson:"age"`                   // 重抓执行周期，单位/秒
}

// 网络请求相关参数
//"method",
//"headers",
//"user_agent",
//"cookies",
//"use_gzip",
//
//"etag",
//"last_modifed",
//"last_modified",
//
//"data",
//
//"proxy",
//"retries",
//"max_redirects",
//
//"connect_timeout",
//"timeout",
//
//"save",
//
//"js_run_at",
//"js_script",
//"js_viewport_width",
//"js_viewport_height",
//
//"load_images",
//"fetch_type",
//
//"validate_cert",
//"robots_txt",
//
type crawlTaskFetch struct {
	// 请求头部信息
	Method    string            `json:"method" bson:"method"`         // http请求方法，默认GET
	Headers   map[string]string `json:"headers" bson:"headers"`       // http header
	UserAgent string            `json:"user_agent" bson:"user_agent"` // userAgent,如果设置了，将会覆盖headers中的User-Agent
	Cookies   map[string]string `json:"cookies" bson:"cookies"`       //请求cookie字典
	UseGzip   bool              `json:"use_gzip" bson:"use_gzip"`

	// TODO 暂时没有考虑，以后再优化
	Etag        string `json:"etag" bson:"etag"`
	LastModifed string `json:"last_modifed" bson:"last_modifed"`

	//请求数据部分
	Data string `json:"data" bson:"data"` // POST请求的数据body，经过url encode编码之后

	//网络控制
	Proxy        string `json:"proxy" bson:"proxy"`                 // http代理, E.g: ip:port
	Retries      int    `json:"retries" bson:"retries"`             // 网络失败时重试次数，默认3次
	MaxRedirects int    `json:"max_redirects" bson:"max_redirects"` //最大重定向次数，默认5次，0不限制，负数不允许重定向

	//超时时间
	ConnectTimeout int `json:"connect_timeout" bson:"connect_timeout"` // 网络连接超时时间，秒
	Timeout        int `json:"timeout" bson:"timeout"`                 // 整体超时时间，秒

	//数据透传
	Save map[string]interface{} `json:"save" bson:"save"`
}

//数据解析处理相关参数
//"callback",
//"process_time_limit",
type crawlTaskProcess struct {
	Callback         string `json:"callback" bson:"callback"`                     // 回调函数对应的key
	ProcessTimeLimit int64  `json:"process_time_limit" bson:"process_time_limit"` // 处理时间限制, 秒
}
