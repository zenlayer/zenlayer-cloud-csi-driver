package common

import (
	"bytes"
	"io"
	"strings"
)

var (
	_ Request = &BaseRequest{}
)

const (
	POST = "POST"
	GET  = "GET"

	HTTP   = "http"
	HTTPS  = "https"
	DOMAIN = "console.zenlayer.com"
)

type Request interface {
	GetScheme() string
	GetHttpMethod() string
	GetDomain() string
	GetAction() string
	GetBodyReader() io.Reader
	GetParams() map[string]string
	GetFormParams() map[string]string
	GetBody() []byte
	GetPath() string
	GetUrl() string
	GetContentType() string
	GetHeaders() map[string]string
	GetApiVersion() string
	GetServiceName() string
	GetMaxAttempts() int
	GetAutoRetries() bool
	GetBodyStr() string
	SetScheme(string)
	SetDomain(string)
	SetPath(string)
	SetContentType(string)
	SetHeaders(headers map[string]string)
	SetServiceName(string)
	SetBody([]byte)
	SetAutoRetries(bool)
	SetMaxAttempts(int)
	//GetReadTimeout() time.Duration
	//GetConnectTimeout() time.Duration
	//GetAutoRetry() bool
	//GetMaxAttempts() int
	//SetReadTimeout(readTimeout time.Duration)
	//SetConnectTimeout(connectTimeout time.Duration)
}

type BaseRequest struct {
	httpMethod  string
	scheme      string
	domain      string
	params      map[string]string
	formParams  map[string]string
	headers     map[string]string
	contentType string
	body        []byte
	path        string

	action      string
	serviceName string
	apiVersion  string
	autoRetries bool `default:"false"`
	maxAttempts int  `default:"0"`
}

func (br *BaseRequest) SetMaxAttempts(maxAttempts int) {
	br.maxAttempts = maxAttempts
}

func (br *BaseRequest) SetAutoRetries(autoRetires bool) {
	br.autoRetries = autoRetires
}

func (br *BaseRequest) GetApiVersion() string {
	return br.apiVersion
}

func (br *BaseRequest) GetBodyStr() string {
	if br.body != nil && len(br.body) > 0 {
		return string(br.body)
	}
	return ""
}

func (br *BaseRequest) SetBody(body []byte) {
	br.body = body
}

func (br *BaseRequest) GetMaxAttempts() int {
	return br.maxAttempts
}

func (br *BaseRequest) GetAutoRetries() bool {
	return br.autoRetries
}

func (br *BaseRequest) SetPath(path string) {
	br.path = path
}

func (br *BaseRequest) GetScheme() string {
	return br.scheme
}

func (br *BaseRequest) GetHttpMethod() string {
	return br.httpMethod
}

func (br *BaseRequest) GetDomain() string {
	return br.domain
}

func (br *BaseRequest) GetAction() string {
	return br.action
}

func (br *BaseRequest) GetBodyReader() io.Reader {
	if br.httpMethod == POST && br.formParams != nil && len(br.formParams) > 0 {
		s := GetUrlQueriesEncoded(br.formParams)
		return strings.NewReader(s)
	} else if br.httpMethod == POST && br.body != nil && len(br.body) > 0 {
		return bytes.NewReader(br.body)
	} else {
		return strings.NewReader("")
	}
}

func (br *BaseRequest) GetParams() map[string]string {
	return br.params
}

func (br *BaseRequest) GetFormParams() map[string]string {
	return br.formParams
}

func (br *BaseRequest) GetBody() []byte {
	return br.body
}

func (br *BaseRequest) GetPath() string {
	return br.path
}

func (br *BaseRequest) GetUrl() string {
	if br.httpMethod == GET {
		return br.GetScheme() + "://" + br.domain + br.path + "?" + GetUrlQueriesEncoded(br.params)
	} else if br.httpMethod == POST {
		return br.GetScheme() + "://" + br.domain + br.path
	} else {
		return ""
	}
}

func (br *BaseRequest) GetContentType() string {
	return br.contentType
}

func (br *BaseRequest) GetHeaders() map[string]string {
	return br.headers
}

func (br *BaseRequest) GetServiceName() string {
	return br.serviceName
}

func (br *BaseRequest) SetScheme(scheme string) {
	scheme = strings.ToLower(scheme)
	switch scheme {
	case HTTP:
		br.scheme = HTTP
	default:
		br.scheme = HTTPS
	}
}

func (br *BaseRequest) SetDomain(domain string) {
	if domain == "" {
		br.domain = DOMAIN
	} else {
		br.domain = domain
	}
}

func (br *BaseRequest) SetContentType(contentType string) {
	br.contentType = contentType
	br.addHeaderParam("Content-Type", contentType)
}

func (br *BaseRequest) addHeaderParam(key, value string) {
	br.headers[key] = value
}

func (br *BaseRequest) SetHeaders(headers map[string]string) {
	if headers == nil {
		return
	}
	br.headers = headers
}

func (br *BaseRequest) SetServiceName(serviceName string) {
	br.serviceName = serviceName
}

func (br *BaseRequest) Init() *BaseRequest {
	br.domain = ""
	br.formParams = make(map[string]string)
	br.params = make(map[string]string)
	br.headers = make(map[string]string)
	br.httpMethod = POST
	return br
}

func (br *BaseRequest) InitWithApiInfo(serviceName, apiVersion, action string) {
	br.serviceName = serviceName
	br.apiVersion = apiVersion
	br.action = action
	br.path = "/api/v2/" + serviceName
}
