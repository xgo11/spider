package fetcher

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

import (
	"github.com/xgo11/datetime"
	"github.com/xgo11/spider/core"
)

const (
	defaultConnectTimeout = 30 * time.Second
	defaultReadTimeout    = 120 * time.Second
	defaultRetryInterval  = 3 * time.Second
	defaultRetryTimes     = 3
	defaultMaxRedirects   = 10
)

type httpClient struct{}

type clientParams struct {
	url            *url.URL
	proxy          *url.URL
	header         http.Header
	method         string
	body           io.Reader
	connectTimeout time.Duration
	readTimeout    time.Duration
	retryInterval  time.Duration
	retryTimes     int
	redirectTimes  int
	client         *http.Client
}

func (hClient *httpClient) Do(task *core.Task) (resp *core.Response) {
	var params *clientParams
	var httpResp *http.Response
	var err error

	resp = &core.Response{Url: task.Url}
	startTime := datetime.Now()

	if params, err = hClient.buildHttpParams(task); err == nil {
		params.client = hClient.buildHttpClient(params)
		if httpResp, err = hClient.doRequest(params); err == nil {
			hClient.deCompressRespBody(httpResp)
		}
	}

	if err == nil && httpResp != nil {
		resp.StatusCode = httpResp.StatusCode
		if currentUrl := httpResp.Request.URL.String(); currentUrl != task.Url {
			resp.Url = currentUrl
			resp.OrigUrl = task.Url
		}
		if len(httpResp.Header) > 0 {
			resp.Headers = httpResp.Header
		}

		resp.Cookies = make(map[string]string)
		if params != nil && params.client.Jar != nil {
			for _, item := range params.client.Jar.Cookies(params.url) {
				resp.Cookies[item.Name] = item.Value
			}
			for _, item := range params.client.Jar.Cookies(httpResp.Request.URL) {
				resp.Cookies[item.Name] = item.Value
			}
		}

		for _, item := range httpResp.Cookies() {
			resp.Cookies[item.Name] = item.Value
		}

		var ctx []byte
		if ctx, err = ioutil.ReadAll(httpResp.Body); err == nil {
			_ = httpResp.Body.Close()
			resp.Content = ctx
			resp.ContentLength = len(resp.Content)
		}
	}
	if err != nil {
		if params == nil { // before request error
			resp.StatusCode = 99
		} else { // http request error
			resp.StatusCode = 599
		}
		resp.ErrMessage = err.Error()
	}
	endTime := datetime.Now()
	resp.TimeMS = int(endTime.Sub(startTime).Round(time.Millisecond))

	task.LastCrawl = endTime.Unix()
	task.Status = core.TaskStatusCrawled

	return
}

func (hClient *httpClient) buildHttpParams(req *core.Task) (param *clientParams, err error) {

	param = &clientParams{
		redirectTimes:  defaultMaxRedirects,
		retryTimes:     defaultRetryTimes,
		retryInterval:  defaultRetryInterval,
		connectTimeout: defaultConnectTimeout,
		readTimeout:    defaultReadTimeout,
		header:         make(http.Header),
	}

	if param.url, err = url.Parse(req.Url); err != nil {
		return nil, err
	}

	var strProxy = req.Fetch.Proxy
	if strProxy != "" {
		if strings.Index(strProxy, "://") < 0 {
			strProxy = "http://" + strProxy
		}
		if param.proxy, err = url.Parse(strProxy); err != nil {
			return nil, err
		}
	}

	var ckMap = make(map[string]string)
	for k, v := range req.Fetch.Headers {
		if strings.ToLower(k) == "cookie" {
			if v != "" {
				for _, s := range strings.Split(v, ";") {
					s = strings.Trim(s, " ")
					if p := strings.Index(s, "="); p > 0 {
						ckMap[strings.Trim(s[0:p], " ")] = strings.Trim(s[p+1:], " ")
					}
				}
				param.header.Add("Cookie", v)
			}
		} else {
			param.header.Add(k, v)
		}
	}

	for k, v := range req.Fetch.Cookies {
		ckMap[k] = v
	}

	if ckSize := len(ckMap); ckSize > 0 {
		var kvArr = make([]string, 0, ckSize)
		for k, v := range ckMap {
			kvArr = append(kvArr, fmt.Sprintf("%s=%s", k, v))
		}
		param.header.Set("Cookie", strings.Join(kvArr, "; "))
	}

	switch method := strings.ToUpper(req.Fetch.Method); method {
	case "GET", "HEAD":
		param.method = method
	case "POST":
		param.method = method
		if contentType := param.header.Get("Content-Type"); len(contentType) < 1 {
			param.header.Add("Content-Type", "application/x-www-form-urlencoded")
		}
		param.body = strings.NewReader(req.Fetch.Data)
	default:
		param.method = "GET"
	}
	if req.Fetch.MaxRedirects > 0 {
		param.redirectTimes = req.Fetch.MaxRedirects
	}

	var connectTimeout time.Duration
	if req.Fetch.ConnectTimeout > 0 {
		connectTimeout = time.Duration(req.Fetch.ConnectTimeout)
	}
	var readTimeout time.Duration
	if req.Fetch.Timeout > 0 {
		readTimeout = time.Duration(req.Fetch.Timeout)
	}

	if connectTimeout > readTimeout {
		connectTimeout = readTimeout
	}
	if connectTimeout > 0 {
		param.connectTimeout = connectTimeout * time.Second
	}
	if readTimeout > 0 {
		param.readTimeout = readTimeout * time.Second
	}

	return
}

func (param *clientParams) checkRedirect(req *http.Request, via []*http.Request) error {

	if param.redirectTimes == 0 {
		return nil
	}
	if len(via) >= param.redirectTimes {
		if param.redirectTimes < 0 {
			return fmt.Errorf("not allow redirects")
		}
		return fmt.Errorf("stopped after %v redirects", param.redirectTimes)
	}

	newCks := req.Cookies()
	oldCks := via[len(via)-1].Cookies()

	if oldCks != nil && len(oldCks) > 0 {
		param.client.Jar.SetCookies(req.URL, oldCks)
	}
	if newCks != nil && len(newCks) > 0 {
		param.client.Jar.SetCookies(req.URL, newCks)
	}
	return nil
}

func (hClient *httpClient) buildHttpClient(param *clientParams) *http.Client {
	client := &http.Client{
		CheckRedirect: param.checkRedirect,
	}

	client.Jar, _ = cookiejar.New(nil)
	var cks []*http.Cookie
	if len(param.header.Get("Cookie")) > 1 {
		for _, part := range strings.Split(param.header.Get("Cookie"), ";") {
			if pos := strings.Index(part, "="); pos > 0 {
				name := strings.Trim(part[0:pos], " ")
				value := strings.Trim(part[pos+1:], " ")
				cks = append(cks, &http.Cookie{Name: name, Value: value})
			}
		}
	}
	if len(cks) > 0 {
		client.Jar.SetCookies(param.url, cks)
	}

	transport := &http.Transport{DialContext: (&net.Dialer{Timeout: param.readTimeout}).DialContext}
	if param.proxy != nil {
		transport.Proxy = http.ProxyURL(param.proxy)
	}
	if strings.ToLower(param.url.Scheme) == "https" {
		transport.TLSClientConfig = &tls.Config{RootCAs: nil, InsecureSkipVerify: true}
		transport.DisableCompression = true
	}
	client.Transport = transport
	return client
}

func (hClient *httpClient) doRequest(param *clientParams) (resp *http.Response, err error) {
	if req, eReq := http.NewRequest(param.method, param.url.String(), param.body); eReq != nil {
		return nil, eReq
	} else {
		req.Header = param.header
		if param.retryTimes <= 0 { // not allowed to retry
			resp, err = param.client.Do(req)
		} else {
			for i := 0; i < param.retryTimes; i++ {
				if resp, err = param.client.Do(req); err == nil {
					break
				}
				if i+1 < param.retryTimes {
					time.Sleep(param.retryInterval)
				}
			}
		}
	}
	return
}

func (hClient *httpClient) deCompressRespBody(resp *http.Response) {
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		var gzipReader *gzip.Reader
		gzipReader, err := gzip.NewReader(resp.Body)
		if err == nil {
			resp.Body = gzipReader
		}

	case "deflate":
		resp.Body = flate.NewReader(resp.Body)

	case "zlib":
		var readCloser io.ReadCloser
		readCloser, err := zlib.NewReader(resp.Body)
		if err == nil {
			resp.Body = readCloser
		}
	}
}
