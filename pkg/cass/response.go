package cass

import (
	"encoding/json"
	"strconv"
)

type Responder interface {
	JSONParse(j string) error
	String() string
	Error() error
	SetError(error)
	HasError() bool
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
	Content     string
	Message     string
	Sign        string
	ResponseKey string
	SubCode     string
	SubMsg      string
	err         error
}

func (resp Response) String() string {
	b, _ := json.Marshal(map[string]string{
		"http":    strconv.Itoa(resp.HTTP.StatusCode),
		"code":    resp.Code,
		"message": resp.Message,
		"subCode": resp.SubCode,
		"subMsg":  resp.SubMsg,
		"content": resp.Content,
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
