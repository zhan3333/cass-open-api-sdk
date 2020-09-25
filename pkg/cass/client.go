package cass

import (
	"encoding/json"
	"go-skysharing-openapi/pkg/signer"
	"io/ioutil"
	"net/http"
)

var F *Factory

type Factory struct {
	Uri             string
	AppId           string
	UserPrivateKey  string
	USigner         signer.Signer
	SSigner         signer.Signer
	UserPublicKey   string
	SystemPublicKey string
	Version         string
	SignType        string
	Charset         string
	Format          string
	HttpClient      http.Client
}

type FactoryConf struct {
	Uri             string
	AppId           string
	UserPublicKey   string
	UserPrivateKey  string
	SystemPublicKey string
}

func NewFactory(c FactoryConf) (*Factory, error) {
	var err error
	var f Factory
	f.Uri = c.Uri
	f.Format = "JSON"
	f.Charset = "UTF-8"
	f.Version = "1.0"
	f.SignType = "RSA2"
	f.AppId = c.AppId
	f.UserPrivateKey = c.UserPrivateKey
	f.UserPublicKey = c.UserPublicKey
	f.SystemPublicKey = c.SystemPublicKey
	f.USigner, err = signer.New(f.UserPrivateKey, f.UserPublicKey)
	if err != nil {
		return &f, err
	}
	f.SSigner, err = signer.New("", f.SystemPublicKey)
	if err != nil {
		return &f, err
	}
	f.HttpClient = *http.DefaultClient
	F = &f
	return &f, nil
}

func (f *Factory) NewRequest(method Method) *Request {
	var request Request
	request.Uri = f.Uri
	request.Params = &params{
		Method:   method.Method,
		AppId:    f.AppId,
		Format:   f.Format,
		Charset:  f.Charset,
		Datetime: "",
		Version:  f.Version,
		SignType: f.SignType,
		BizParam: "",
		Sign:     "",
	}
	request.signer = f.USigner
	request.HttpClient = f.HttpClient
	return &request
}

func (f *Factory) NewResponse(resp *http.Response, responseKey string) *Response {
	response := new(Response)
	response.Key = responseKey
	response.StatusCode = resp.StatusCode
	response.Body, _ = ioutil.ReadAll(resp.Body)
	response.Json = string(response.Body)
	_ = json.Unmarshal(response.Body, &response.Map)
	if _, ok := response.Map[responseKey]; !ok {
		return response
	}
	response.Code = response.Map[responseKey].(map[string]interface{})["code"].(string)
	response.Content = response.Map[responseKey].(map[string]interface{})["content"].(string)
	response.Message = response.Map[responseKey].(map[string]interface{})["message"].(string)
	response.Sign = response.Map["sign"].(string)
	return response
}
