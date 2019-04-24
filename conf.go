package spider

type Config struct {
	MessageQueue string `yaml:"message_queue"`
	Scheduler    struct {
		HTTPServer string `yaml:"http_server"`
	} `yaml:"scheduler"`
	Fetcher struct {
		HTTPServer string `yaml:"http_server"`
	} `yaml:"fetcher"`

	Processor struct {
		HTTPServer string `yaml:"http_server"`
	} `yaml:"processor"`

	ResultWorker struct {
	} `yaml:"result_worker"`

	Web struct {
		HTTPServer string `yaml:"http_server"`
	} `yaml:"web"`
}
