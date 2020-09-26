package cass

import (
	"encoding/json"
	"strconv"
)

type Response struct {
	StatusCode int
	Body       []byte
	Json       string
	Map        map[string]interface{}
	Code       string
	Content    string
	Message    string
	Sign       string
	Key        string
	SubCode    string
	SubMsg     string
}

func (resp Response) String() string {
	b, _ := json.Marshal(map[string]string{
		"http":    strconv.Itoa(resp.StatusCode),
		"code":    resp.Code,
		"message": resp.Message,
		"subCode": resp.SubCode,
		"subMsg":  resp.SubMsg,
		"content": resp.Content,
	})
	return string(b)
}
