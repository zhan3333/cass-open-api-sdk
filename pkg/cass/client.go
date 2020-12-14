package cass

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-skysharing-openapi/pkg/cass/method"
	"go-skysharing-openapi/pkg/signer"
	"io/ioutil"
	"net/http"
	"strings"
)

var F *Client

type Client struct {
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
	Debug           bool
}

type Config struct {
	URI             string
	AppId           string
	UserPublicKey   string
	UserPrivateKey  string
	SystemPublicKey string
	Debug           bool
}

func NewClient(c Config) (*Client, error) {
	var err error
	var f Client
	f.Debug = c.Debug
	f.Uri = c.URI
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

func (f *Client) NewRequest(method method.Method) Requester {
	var request Request
	request.Debug = f.Debug
	request.Uri = f.Uri
	request.Params = &Params{
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
	request.Method = method
	return &request
}

func NewResponse(resp *http.Response, responseKey string, err error) Responder {
	var response Response
	response.SetError(err)
	if response.HasError() {
		return &response
	}
	response.ResponseKey = responseKey
	response.HTTP.StatusCode = resp.StatusCode
	response.HTTP.Body, _ = ioutil.ReadAll(resp.Body)
	response.HTTP.JSON = string(response.HTTP.Body)
	err = json.Unmarshal(response.HTTP.Body, &response.HTTP.Params)
	if err != nil {
		response.SetError(err)
		return &response
	}
	if _, ok := response.HTTP.Params[responseKey]; !ok {
		response.SetError(errors.New(fmt.Sprintf("response not have %s key", responseKey)))
		return &response
	}
	response.Code = response.HTTP.Params[responseKey].(map[string]interface{})["code"].(string)
	response.ContentStr = response.HTTP.Params[responseKey].(map[string]interface{})["content"].(string)
	response.Message = response.HTTP.Params[responseKey].(map[string]interface{})["message"].(string)
	response.SubCode = response.HTTP.Params[responseKey].(map[string]interface{})["subCode"].(string)
	response.SubMsg = response.HTTP.Params[responseKey].(map[string]interface{})["subMsg"].(string)
	response.Sign = response.HTTP.Params["sign"].(string)
	return &response
}

func (f *Client) Send(request Requester) Responder {
	var err error
	var response *http.Response
	if f.Debug {
		fmt.Printf("request: %+v\n", request)
	}
	query, err := request.GetQuery()
	if f.Debug {
		fmt.Printf("query: %+v; err: %+v\n", query, err)
	}
	if err != nil {
		return NewResponse(nil, request.GetMethod().GetResponseKey(), err)
	}
	response, err = f.HttpClient.Post(
		f.Uri,
		"application/html; charset=utf-8",
		strings.NewReader(query),
	)
	if f.Debug {
		fmt.Printf("response: %+v; err: %+v\n", response, err)
	}
	return NewResponse(response, request.GetMethod().GetResponseKey(), err)
}
