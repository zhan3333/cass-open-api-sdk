package cass

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-skysharing-openapi/pkg/cass/context"
	"go-skysharing-openapi/pkg/cass/method"
	"go-skysharing-openapi/pkg/signer"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Requester interface {
	// 设置业务参数
	SetBizParams(b context.Biz)
	// 获取查询参数
	GetQuery() (string, error)
	GetMethod() method.Method
}

type Request struct {
	Requester
	Uri        string
	signer     signer.Signer
	BizParam   context.Biz
	SignStr    string `json:"sign"`
	Params     *Params
	HttpClient http.Client
	Debug      bool
	Method     method.Method
}

type Params struct {
	Method   string `json:"method"`
	AppId    string `json:"APPID"`
	Format   string `json:"format"`
	Charset  string `json:"charset"`
	Datetime string `json:"datetime"`
	Version  string `json:"version"`
	SignType string `json:"signType"`
	BizParam string `json:"bizParam"`
	Sign     string `json:"sign"`
}

func (params *Params) toMap() (map[string]interface{}, error) {
	var waitSignParams = map[string]interface{}{}

	j, _ := json.Marshal(params)
	err := json.Unmarshal(j, &waitSignParams)
	if err != nil {
		return waitSignParams, err
	}
	return waitSignParams, nil
}

func (request *Request) SetBizParams(b context.Biz) {
	request.BizParam = b
}

func sign(request *Request) (string, error) {
	var err error
	bizParamBytes, _ := json.Marshal(request.BizParam)
	request.Params.BizParam = string(bizParamBytes)
	request.Params.Datetime = time.Now().Format("2006-01-02 15:04:05")
	// 将 bizParam 转为json, 其中的中文不要转为 unicode 编码, 保持中文字符
	waitSignParams, err := request.Params.toMap()
	if err != nil {
		return "", err
	}

	// 过滤 Request 中的空字符: '', null, '[]', '{}'
	for s, v := range waitSignParams {
		if v == "" || v == "{}" || v == "[]" || v == nil {
			delete(waitSignParams, s)
		}
	}

	// 将 key 按照升序排序
	sortedKeys := make([]string, 0)
	for k := range waitSignParams {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	signMapParams := make(map[string]interface{}, 0) //加密使用
	for _, k := range sortedKeys {
		v := waitSignParams[k]
		if k != "" && v != "" {
			if v == "" || v == "{}" || v == "[]" {
				continue
			}
			signMapParams[k] = v
		}
	}

	// 将 request 转换为 json
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(signMapParams)
	if err != nil {
		return "", err
	}
	jsonStr := bf.String()
	// 将 request json str 中的空格 (ASCII 码空格) 去掉
	// 服务端的 json_encode 会将 // 转义为 \/\/, 但是 golang 中不会转义
	jsonStr = strings.ReplaceAll(jsonStr, " ", "")
	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
	jsonStr = strings.ReplaceAll(jsonStr, "/", `\/`)

	// 将 request json str 进行 urlencode 编码, 产生待签名字符串
	urlEncodeStr := url.QueryEscape(jsonStr)
	// 通过字符串生成签名
	signBytes, err := request.signer.Sign([]byte(urlEncodeStr), crypto.SHA256)
	if err != nil {
		return "", err
	}
	sign := base64.StdEncoding.EncodeToString(signBytes)

	// 对 sign 进行 urlencode 处理, 防止 base64 中的字符串在 url 无法正常作为参数
	return sign, nil
}

func (request *Request) GetQuery() (string, error) {
	sign, err := sign(request)
	if err != nil {
		return "", err
	}
	request.SignStr = sign
	request.Params.Sign = sign

	// 生成请求字符串
	str := new(strings.Builder)
	waitBuildQueryParams, err := request.Params.toMap()
	if err != nil {
		return "", err
	}
	if len(waitBuildQueryParams) != 0 {
		for key, val := range waitBuildQueryParams {
			str.WriteString(fmt.Sprintf("%s=%s&", key, url.QueryEscape(fmt.Sprintf("%s", val))))
		}
	}
	return strings.TrimRight(str.String(), "&"), nil
}

func (request *Request) GetMethod() method.Method {
	return request.Method
}
