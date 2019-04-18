package model

import "fmt"

var (
	ErrorSessionTimeout  = fmt.Errorf("session_timeout")
	ErrorPipelineTimeout = fmt.Errorf("pipeline_timeout")
)

//定义解析结果
type ProcessResult struct {
	// 解析数据
	Error      string      `json:"error"`       // 解析错误
	Html       string      `json:"html"`        // 页面原始信息
	StatusCode int         `json:"status_code"` // http状态码
	Result     interface{} `json:"result"`      // 解析后的结构化数据
	Url        string      `json:"url"`
	OrigUrl    string      `json:"orig_url"`

	// 任务信息
	TaskId  string                 `json:"task_id"`  // 任务id
	Project string                 `json:"project"`  // 所属项目
	Catg    string                 `json:"catg"`     // 类别
	SubCatg string                 `json:"sub_catg"` // 二级类别
	AuthId  string                 `json:"auth_id"`  // 授权id
	Save    map[string]interface{} `json:"save"`     // 任务透传数据

	// 附加
	NeedSnapshot bool `json:"need_snapshot"` // 是否需要保存快照，某些情况下需要快照分析页面
}

func NewProcessResult(result *FetchResult) *ProcessResult {
	pr := &ProcessResult{}
	pr.Init(result)
	return pr
}
func (pr *ProcessResult) Init(result *FetchResult) {
	if len(pr.Error) < 1 {
		pr.Error = result.Error
	}
	pr.Html = result.GetText()
	pr.StatusCode = result.StatusCode
	pr.Url = result.Url
	pr.OrigUrl = result.OrigUrl
	pr.TaskId = result.Task.TaskId
	pr.Catg = result.Task.Catg
	pr.SubCatg = result.Task.SubCatg
	pr.Save = result.Task.GetSave()
	pr.Project = result.Task.Project
	pr.AuthId = result.Task.GetAuthId()
}

func (pr *ProcessResult) HasError() bool {

	if pr.Error != "" {
		return true
	}
	return false

}
