package cass_test

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go-skysharing-openapi/pkg/cass"
	"net/url"
	"os"
	"testing"
)

var factoryConf cass.FactoryConf

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env")
	factoryConf = cass.FactoryConf{
		Uri:             os.Getenv("API_URL"),
		AppId:           os.Getenv("APPID"),
		UserPublicKey:   os.Getenv("PUBLIC_KEY_STR"),
		UserPrivateKey:  os.Getenv("PRIVATE_KEY_STR"),
		SystemPublicKey: os.Getenv("VZHUO_PUBLIC_KEY_STR"),
	}
	m.Run()
}

func TestGetBalance(t *testing.T) {
	f, err := cass.NewFactory(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(cass.M.GetBalance)
	request.BizParam = map[string]interface{}{}
	response, err := request.Send()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.StatusCode)
	t.Log(response.String())
}

func BenchmarkOneBankPayParallel(b *testing.B) {
	f, err := cass.NewFactory(factoryConf)
	assert.Nil(b, err)
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			request := f.NewRequest(cass.M.PayOneBankRemit)
			request.BizParam = map[string]interface{}{
				"payChannelK": "1",
				//"payeeChannelType": "2",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "13517210601",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			}
			response, err := request.Send()
			assert.Nil(b, err)
			assert.NotNil(b, response)
			assert.Equal(b, 200, response.StatusCode)
			b.Log(response.String())
		}
	})
}

func TestOneBankPay(t *testing.T) {
	f, err := cass.NewFactory(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(cass.M.PayOneBankRemit)
	request.BizParam = map[string]interface{}{
		"payChannelK": "1",
		//"payeeChannelType": "2",
		"orderData": [1]interface{}{
			map[string]interface{}{
				"orderSN":          uuid.New().String(),
				"receiptFANO":      "13517210601",
				"payeeAccount":     "詹光",
				"requestPayAmount": "0.01",
				"notifyUrl":        "http://www.baidu.com/a/b?a=b",
			},
		},
	}
	response, err := request.Send()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.StatusCode)
	t.Log(response.String())
}

func TestMaYiBankPay(t *testing.T) {
	f, err := cass.NewFactory(factoryConf)
	assert.Nil(t, err)

	var tests = []struct {
		params     map[string]interface{}
		expectResp cass.Response
	}{
		// 无收款通道
		{
			map[string]interface{}{
				"payChannelK": "1",
				//"payeeChannelType": "2",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "13517210601",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				StatusCode: 200,
				Code:       "50000",
				Message:    "业务参数校验错误",
				SubCode:    "api.common.error",
				SubMsg:     "收款通道 不能为空。",
			},
		},
		// 收款通道错误
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "3",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "13517210601",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				StatusCode: 200,
				Code:       "50000",
				Message:    "业务参数校验错误",
				SubCode:    "api.common.error",
				SubMsg:     "已选的属性 收款通道 非法。",
			},
		},
		// 错误的收款账号
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "1", // 银行卡
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "abcd",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				StatusCode: 200,
				Code:       "50000",
				Message:    "业务参数校验错误",
				SubCode:    "api.common.error",
				SubMsg:     "收款人账号格式错误，不是一个有效的银行卡号或者银行户口号。",
			},
		},
	}

	for _, tt := range tests {
		request := f.NewRequest(cass.M.PayOneBankRemit)
		request.BizParam = tt.params
		response, err := request.Send()
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, tt.expectResp.StatusCode, response.StatusCode)
		assert.Equal(t, tt.expectResp.Code, response.Code)
		assert.Equal(t, tt.expectResp.SubCode, response.SubCode)
		assert.Equal(t, tt.expectResp.Message, response.Message)
		assert.Equal(t, tt.expectResp.SubMsg, response.SubMsg)
		t.Log(response.String())
	}
}

func TestMaYiBankPayNotAllowPayeeChannelType(t *testing.T) {
	f, err := cass.NewFactory(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(cass.M.PayOneBankRemit)
	request.BizParam = map[string]interface{}{
		"payChannelK":      "1",
		"payeeChannelType": "3",
		"orderData": [1]interface{}{
			map[string]interface{}{
				"orderSN":          uuid.New().String(),
				"receiptFANO":      "13517210601",
				"payeeAccount":     "詹光",
				"requestPayAmount": "0.01",
				"notifyUrl":        "http://www.baidu.com/a/b?a=b",
			},
		},
	}
	response, err := request.Send()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "已选的属性 收款通道 非法。", response.SubMsg)
	t.Log(response.String())
}

func TestOneWeChatPay(t *testing.T) {
	f, err := cass.NewFactory(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(cass.M.PayOneWeChatRemit)
	request.SetBizParams(cass.PayOneWeChatRemitBiz{
		PayChannelK: "3",
		OrderData: []cass.WeChatOrder{
			{
				OrderSN:          uuid.New().String(),
				Phone:            "13517210601",
				PayeeAccount:     "詹光",
				RequestPayAmount: "1",
				NotifyUrl:        "http://www.baidu.com/",
			},
		},
	})
	response, err := request.Send()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.StatusCode)
	t.Log(response.String())
}

func TestUrlQueryEscape(t *testing.T) {
	s := url.QueryEscape("http://www.baidu.com")
	t.Log(s)
	bites, _ := json.Marshal(map[string]string{
		"url": "http://www.baidu.com?name=zhan&age=22",
	})
	t.Logf("%s", bites)

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	_ = jsonEncoder.Encode(map[string]string{
		"url": "http://www.baidu.com?name=zhan&age=22",
	})
	t.Log(bf.String())
}

func TestFactory_New(t *testing.T) {
	var err error
	cass.F, err = cass.NewFactory(factoryConf)
	assert.Nil(t, err)
	assert.NotNil(t, cass.F)
}

func TestFactory_NewRequest(t *testing.T) {
	var err error
	cass.F, err = cass.NewFactory(factoryConf)
	assert.Nil(t, err)
	assert.NotNil(t, cass.F)

	request := cass.F.NewRequest(cass.M.GetBalance)
	request.BuildParams()
}
