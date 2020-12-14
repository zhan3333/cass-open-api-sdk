package cass

import (
	"encoding/json"
	"go-skysharing-openapi/pkg/cass/context"
	"strconv"
)

type Responder interface {
	JSONParse(j string) error
	String() string
	Error() error
	SetError(error)
	HasError() bool
	Content(context.Content) error
	// 请求是否成功
	IsHTTPSuccess() bool
	// 操作是否成功
	IsSuccess() bool
	// 业务是否处理成功
	IsBusinessSuccess() bool
}

type ResponseHTTP struct {
	StatusCode int
	Body       []byte
	JSON       string
	Params     map[string]interface{}
}

type Response struct {
	Responder
	HTTP ResponseHTTP
	Code string
	// 业务响应字符串
	Message     string
	Sign        string
	ResponseKey string
	SubCode     string
	SubMsg      string
	ContentStr  string
	err         error
}

func (resp Response) String() string {
	b, _ := json.Marshal(map[string]string{
		"http":    strconv.Itoa(resp.HTTP.StatusCode),
		"code":    resp.Code,
		"message": resp.Message,
		"subCode": resp.SubCode,
		"subMsg":  resp.SubMsg,
		"content": resp.ContentStr,
	})
	return string(b)
}

func (resp *Response) JSONParse(j string) error {
	return nil
}

func (resp *Response) Error() error {
	return resp.err
}
func (resp *Response) SetError(err error) {
	resp.err = err
}
func (resp *Response) HasError() bool {
	return resp.err != nil
}
func (resp *Response) Content(c context.Content) error {
	return json.Unmarshal([]byte(resp.ContentStr), &c)
}
func (resp *Response) IsHTTPSuccess() bool {
	return resp.HTTP.StatusCode == 200
}
func (resp *Response) IsBusinessSuccess() bool {
	return resp.Code == "10000"
}
func (resp *Response) IsSuccess() bool {
	return resp.IsHTTPSuccess() && resp.IsBusinessSuccess()
}
