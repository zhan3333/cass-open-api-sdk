package cass

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-skysharing-openapi/pkg/signer"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Request struct {
	Uri        string
	signer     signer.Signer
	BizParam   Biz
	Sign       string `json:"sign"`
	Params     *params
	HttpClient http.Client
}

type params struct {
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

func (params *params) toMap() (map[string]interface{}, error) {
	var waitSignParams = map[string]interface{}{}

	j, _ := json.Marshal(params)
	err := json.Unmarshal(j, &waitSignParams)
	if err != nil {
		return waitSignParams, err
	}
	return waitSignParams, nil
}

func (request *Request) SetBizParams(b Biz) {
	request.BizParam = b
}

func (params params) BuildQuery() (string, error) {
	str := new(strings.Builder)
	waitBuildQueryParams, err := params.toMap()
	if err != nil {
		return str.String(), err
	}
	if len(waitBuildQueryParams) != 0 {
		for key, val := range waitBuildQueryParams {
			str.WriteString(fmt.Sprintf("%s=%s&", key, url.QueryEscape(fmt.Sprintf("%s", val))))
		}
	}
	return strings.TrimRight(str.String(), "&"), nil
}

// 构建请求参数对象
func (request *Request) BuildParams() {
	bizParamBytes, _ := json.Marshal(request.BizParam)
	request.Params.BizParam = string(bizParamBytes)
	request.Params.Datetime = time.Now().Format("2006-01-02 15:04:05")
	_ = request.makeSign()
}

func (request *Request) makeSign() error {
	var err error
	// 将 bizParam 转为json, 其中的中文不要转为 unicode 编码, 保持中文字符
	waitSignParams, err := request.Params.toMap()
	if err != nil {
		return err
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
		return err
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
	sign := base64.StdEncoding.EncodeToString(signBytes)

	// 对 sign 进行 urlencode 处理, 防止 base64 中的字符串在 url 无法正常作为参数

	//sign = url.QueryEscape(sign)
	if err != nil {
		return err
	}
	waitSignParams["sign"] = sign
	request.Sign = sign
	request.Params.Sign = sign
	return nil
}

func (request *Request) Send() (*Response, error) {
	var err error
	var response *http.Response
	request.BuildParams()
	query, err := request.Params.BuildQuery()
	if err != nil {
		return nil, err
	}
	response, err = request.HttpClient.Post(
		request.Uri,
		"application/html; charset=utf-8",
		strings.NewReader(query),
	)

	if err != nil {
		return nil, err
	}
	return F.NewResponse(response, request.GetResponseKey()), nil
}

func (request *Request) GetResponseKey() string {
	s := strings.ReplaceAll(request.Params.Method, ".", "")
	return fmt.Sprintf("%sResponse", s)
}

type Biz interface {
}

type PayOneBankRemitBiz struct {
	PayChannelK      string      `json:"payChannelK" binding:"required,max=1" comment:" 付款通道：1-银行卡"`
	PayeeChannelType string      `json:"payeeChannelType" binding:"max=1" comment:"收款通道： 网商银行为必填，1-银行卡，2-支付宝； 其他通道，收款与付款通道一致，不需要接 口传入数据。"`
	ContractID       string      `json:"contractID,omitempty" binding:"" comment:"合同 ID"`
	OrderData        []BankOrder `json:"orderData" binding:"required" comment:" 订单数据，二维数组"`
}

type PayBankRemitBiz struct {
	PayChannelK      string      `json:"payChannelK" binding:"required,max=1" comment:" 付款通道：1-银行卡"`
	PayeeChannelType string      `json:"payeeChannelType" binding:"max=1" comment:"收款通道： 网商银行为必填，1-银行卡，2-支付宝； 其他通道，收款与付款通道一致，不需要接 口传入数据。"`
	ContractID       string      `json:"contractID,omitempty" binding:"" comment:"合同 ID"`
	OrderData        []BankOrder `json:"orderData" binding:"required" comment:" 订单数据，二维数组"`
}

type BankOrder struct {
	OrderSN          string `json:"orderSN" binding:"required,max=64" comment:"商户订单号，只能是英文字母，数字，中文 以及连接符-" example:"测试商户-20190712-0026"`
	ReceiptFANO      string `json:"receiptFANO" binding:"required,max=23" comment:"收款人金融账号： 网商银行支持，1-银行账号 2-支付宝账号，邮箱，手机号 其他银行仅支持，1-银行账号" example:"银行卡：6214832719548955 支付宝：18800080008"`
	PayeeAccount     string `json:"payeeAccount" binding:"required,max=64" comment:"  收款人户名（真实姓名）" example:"张三"`
	ReceiptBankName  string `json:"receiptBankName,omitempty" binding:"max=64" comment:"收款方开户行名称（总行即可）" example:"招商银行"`
	ReceiptBankAddr  string `json:"receiptBankAddr,omitempty" binding:"max=100" comment:"  收款方开户行地（省会或者直辖市即可）" example:"武汉"`
	CRCHGNO          string `json:"CRCHGNO,omitempty" binding:"max=12" comment:"收方联行号" example:" 411234567222"`
	RequestPayAmount string `json:"requestPayAmount" binding:"required,max=26" comment:"预期付款金额" example:"14.00"`
	IdentityCard     string `json:"identityCard,omitempty" binding:"max=20" comment:"收款人身份证号" example:" 321123456789098765"`
	NotifyUrl        string `json:"notifyUrl" binding:"required,max=255" comment:"网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单”" example:" └notifyUrl String -- 255 网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单” http://xxx.xxx.cn/xx/asynNotify.h tm"`
}

type PayOneWeChatRemitBiz struct {
	PayChannelK string        `json:"payChannelK" binding:"required,max=1" comment:" 付款通道：3-微信"`
	ContractID  string        `json:"contractID,omitempty" binding:"" comment:"合同 ID"`
	OrderData   []WeChatOrder `json:"orderData" binding:"required" comment:" 订单数据，二维数组"`
}

type PayWeChatRemitBiz struct {
	PayChannelK string        `json:"payChannelK" binding:"required,max=1" comment:" 付款通道：3-微信"`
	ContractID  string        `json:"contractID,omitempty" binding:"" comment:"合同 ID"`
	OrderData   []WeChatOrder `json:"orderData" binding:"required" comment:" 订单数据，二维数组"`
}

type WeChatOrder struct {
	OrderSN          string `json:"orderSN" binding:"required,max=64" comment:"商户订单号，只能是英文字母，数字，中文 以及连接符-" example:"测试商户-20190712-0026"`
	Phone            string `json:"phone" binding:"required,max=12" comment:"收款人手机号" example:"13517588888"`
	PayeeAccount     string `json:"payeeAccount" binding:"required,max=64" comment:"  收款人户名（真实姓名）" example:"张三"`
	RequestPayAmount string `json:"requestPayAmount" binding:"required,max=26" comment:"预期付款金额" example:"14.00"`
	IdentityCard     string `json:"identityCard,omitempty" binding:"max=20" comment:"收款人身份证号" example:" 321123456789098765"`
	NotifyUrl        string `json:"notifyUrl" binding:"required,max=255" comment:"网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单”" example:" └notifyUrl String -- 255 网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单” http://xxx.xxx.cn/xx/asynNotify.h tm"`
}
