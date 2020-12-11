package go_openapi_test

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go-skysharing-openapi/pkg/cass"
	"go-skysharing-openapi/pkg/cass/context"
	"go-skysharing-openapi/pkg/cass/method"
	"os"
	"strconv"
	"testing"
)

var factoryConf cass.Config

func TestMain(m *testing.M) {
	_ = godotenv.Load(".env")
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	factoryConf = cass.Config{
		URI:             os.Getenv("API_URL"),
		AppId:           os.Getenv("APPID"),
		UserPublicKey:   os.Getenv("PUBLIC_KEY_STR"),
		UserPrivateKey:  os.Getenv("PRIVATE_KEY_STR"),
		SystemPublicKey: os.Getenv("VZHUO_PUBLIC_KEY_STR"),
		Debug:           debug,
	}
	m.Run()
}

// 测试获取余额
func TestGetBalance(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(context.GetBalance)
	request.SetBizParams(context.GetBalanceBiz{
		PayChannelK: "1",
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, response.Error())
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTP.StatusCode)
	assert.Equal(t, "10000", response.Code)
	assert.Equal(t, "请求成功", response.Message)
	assert.Equal(t, "", response.SubCode)
	assert.Equal(t, "", response.SubMsg)
	t.Log(response.String())
}

// 压测付款
func BenchmarkOneBankPayParallel(b *testing.B) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(b, err)
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			request := f.NewRequest(method.M.PayOneBankRemit)
			request.SetBizParams(map[string]interface{}{
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
			})
			response := f.Send(request).(*cass.Response)
			assert.Nil(b, err)
			assert.NotNil(b, response)
			assert.Equal(b, 200, response.HTTP.StatusCode)
			b.Log(response.String())
		}
	})
}

// 测试单笔银行卡付款
func TestOneBankPay(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(method.M.PayOneBankRemit)
	request.SetBizParams(map[string]interface{}{
		"payChannelK":      "1",
		"payeeChannelType": "1",
		"orderData": [1]interface{}{
			map[string]interface{}{
				"orderSN":          uuid.New().String(),
				"receiptFANO":      "6214850278508756",
				"payeeAccount":     "詹光",
				"requestPayAmount": "0.01",
				"notifyUrl":        "http://www.baidu.com/a/b?a=b",
				"data":             "",
			},
		},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTP.StatusCode)
	t.Log(response.String())
}

func TestVerifyIDCardFailed(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(method.M.PayOneBankRemit)
	request.SetBizParams(map[string]interface{}{
		"payChannelK":      "1",
		"payeeChannelType": "1",
		"orderData": [1]interface{}{
			map[string]interface{}{
				"orderSN":          uuid.New().String(),
				"receiptFANO":      "6214850278508756",
				"payeeAccount":     "13517210601",
				"identityCard":     "420222199212041057",
				"requestPayAmount": "0.01",
				"notifyUrl":        "http://www.baidu.com/a/b?a=b",
				"data":             "",
			},
		},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTP.StatusCode)
	t.Log(response.String())
}

// 测试单笔银行卡支付
func TestMaYiBankPayOneWithParamsFailed(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)

	var tests = []struct {
		params     map[string]interface{}
		expectResp cass.Response
	}{
		// 无支付通道
		{
			map[string]interface{}{
				//"payChannelK": "1",
				"payeeChannelType": "2",
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "付款通道不能为空",
			},
		},
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款通道不能为空",
			},
		},
		// 无订单数据
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "2",
				//"orderData": [1]interface{}{
				//	map[string]interface{}{
				//		"orderSN":          uuid.New().String(),
				//		"receiptFANO":      "13517210601",
				//		"payeeAccount":     "詹光",
				//		"requestPayAmount": "0.01",
				//		"notifyUrl":        "http://www.baidu.com/a/b?a=b",
				//	},
				//},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "订单数据不能为空",
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "已选的属性收款通道非法",
			},
		},
		// 错误的收款账号 (错误的银行卡)
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款人账号格式错误，不是一个有效的银行卡号或者银行户口号",
			},
		},
		// 错误的收款账号(错误的银行卡2)
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "1", // 银行卡
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "390961827@qq.com",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款人账号请检查您填写的数字串中是否含有我们不允许的点号(含全角半角字符)",
			},
		},
		// 银行卡成功
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "1", // 银行卡
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "6214850278508756",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
		// 支付宝成功(手机号)
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "2", // 支付宝
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
		// 支付宝成功(邮箱)
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "2", // 支付宝
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "390961827@qq.com",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
	}

	for _, tt := range tests {
		request := f.NewRequest(method.M.PayOneBankRemit)
		request.SetBizParams(tt.params)
		response := f.Send(request).(*cass.Response)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, tt.expectResp.HTTP.StatusCode, response.HTTP.StatusCode)
		assert.Equal(t, tt.expectResp.Code, response.Code)
		assert.Equal(t, tt.expectResp.SubCode, response.SubCode)
		assert.Equal(t, tt.expectResp.Message, response.Message)
		assert.Equal(t, tt.expectResp.SubMsg, response.SubMsg)
		t.Log(response.String())
	}
}

// 测试批量提交蚂蚁银行支付
func TestMaYiBanPayBatch(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)

	var tests = []struct {
		params     map[string]interface{}
		expectResp cass.Response
	}{
		// 无支付通道
		{
			map[string]interface{}{
				//"payChannelK": "1",
				"payeeChannelType": "2",
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "付款通道不能为空",
			},
		},
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款通道不能为空",
			},
		},
		// 无订单数据
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "2",
				//"orderData": [1]interface{}{
				//	map[string]interface{}{
				//		"orderSN":          uuid.New().String(),
				//		"receiptFANO":      "13517210601",
				//		"payeeAccount":     "詹光",
				//		"requestPayAmount": "0.01",
				//		"notifyUrl":        "http://www.baidu.com/a/b?a=b",
				//	},
				//},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "订单数据不能为空",
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "已选的属性收款通道非法",
			},
		},
		// 错误的收款账号 (错误的银行卡)
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款人账号格式错误，不是一个有效的银行卡号或者银行户口号",
			},
		},
		// 错误的收款账号(错误的银行卡2)
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "1", // 银行卡
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "390961827@qq.com",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款人账号请检查您填写的数字串中是否含有我们不允许的点号(含全角半角字符)",
			},
		},
		// 银行卡成功
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "1", // 银行卡
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "6214850278508756",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
		// 支付宝成功(手机号)
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "2", // 支付宝
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
		// 支付宝成功(邮箱)
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "2", // 支付宝
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "390961827@qq.com",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
		// 多笔订单成功
		{
			map[string]interface{}{
				"payChannelK":      "1",
				"payeeChannelType": "2", // 支付宝
				"orderData": [2]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "390961827@qq.com",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"receiptFANO":      "390961827@qq.com",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
	}

	for _, tt := range tests {
		request := f.NewRequest(method.M.PayBankRemit)
		request.SetBizParams(tt.params)
		response := f.Send(request).(*cass.Response)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, tt.expectResp.HTTP.StatusCode, response.HTTP.StatusCode)
		assert.Equal(t, tt.expectResp.Code, response.Code)
		assert.Equal(t, tt.expectResp.SubCode, response.SubCode)
		assert.Equal(t, tt.expectResp.Message, response.Message)
		assert.Equal(t, tt.expectResp.SubMsg, response.SubMsg)
		t.Log(response.String())
	}
}

// 测试微信单笔支付
func TestWeChatPayOne(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)

	var tests = []struct {
		params     map[string]interface{}
		expectResp cass.Response
	}{
		// 无支付通道
		{
			map[string]interface{}{
				//"payChannelK": "1",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"identityCard":     "420222199212041057",
						"phone":            "13517210601",
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "付款通道不能为空",
			},
		},
		// 错误的支付通道
		{
			map[string]interface{}{
				"payChannelK": "0",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"identityCard":     "420222199212041057",
						"phone":            "13517210601",
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "已选的属性付款通道非法",
			},
		},
		// 无订单数据
		{
			map[string]interface{}{
				"payChannelK": "3",
				//"orderData": [1]interface{}{
				//	map[string]interface{}{
				//		"orderSN":          uuid.New().String(),
				//		"receiptFANO":      "13517210601",
				//		"payeeAccount":     "詹光",
				//		"requestPayAmount": "0.01",
				//		"notifyUrl":        "http://www.baidu.com/a/b?a=b",
				//	},
				//},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "订单数据不能为空",
			},
		},
		// 支付通道错误
		{
			map[string]interface{}{
				"payChannelK": "0",
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "已选的属性付款通道非法",
			},
		},
		// 手机号为空
		// 错误的收款账号
		{
			map[string]interface{}{
				"payChannelK": "3",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款人手机号不能为空",
			},
		},
		// 错误的收款账号
		{
			map[string]interface{}{
				"payChannelK": "3",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"phone":            "abcd",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款人手机号不是一个有效的中国内地手机号码",
			},
		},
		// 请求金额需要大于1元
		{
			map[string]interface{}{
				"payChannelK": "3",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"identityCard":     "420222199212041057",
						"phone":            "13517210601",
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "付款金额金额不能少于1元",
			},
		},
		// 成功
		{
			map[string]interface{}{
				"payChannelK": "3",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"identityCard":     "420222199212041057",
						"phone":            "13517210601",
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "1",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
	}

	for _, tt := range tests {
		request := f.NewRequest(method.M.PayOneWeChatRemit)
		request.SetBizParams(tt.params)
		response := f.Send(request).(*cass.Response)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, tt.expectResp.HTTP.StatusCode, response.HTTP.StatusCode)
		assert.Equal(t, tt.expectResp.Code, response.Code)
		assert.Equal(t, tt.expectResp.SubCode, response.SubCode)
		assert.Equal(t, tt.expectResp.Message, response.Message)
		assert.Equal(t, tt.expectResp.SubMsg, response.SubMsg)
		t.Log(response.String())
	}
}

// 测试微信批量支付
func TestWeChatPayBatch(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)

	var tests = []struct {
		params     map[string]interface{}
		expectResp cass.Response
	}{
		// 无支付通道
		{
			map[string]interface{}{
				//"payChannelK": "1",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"identityCard":     "420222199212041057",
						"phone":            "13517210601",
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "付款通道不能为空",
			},
		},
		// 错误的支付通道
		{
			map[string]interface{}{
				"payChannelK": "0",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"identityCard":     "420222199212041057",
						"phone":            "13517210601",
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "已选的属性付款通道非法",
			},
		},
		// 无订单数据
		{
			map[string]interface{}{
				"payChannelK": "3",
				//"orderData": [1]interface{}{
				//	map[string]interface{}{
				//		"orderSN":          uuid.New().String(),
				//		"receiptFANO":      "13517210601",
				//		"payeeAccount":     "詹光",
				//		"requestPayAmount": "0.01",
				//		"notifyUrl":        "http://www.baidu.com/a/b?a=b",
				//	},
				//},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "订单数据不能为空",
			},
		},
		// 支付通道错误
		{
			map[string]interface{}{
				"payChannelK": "0",
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
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "已选的属性付款通道非法",
			},
		},
		// 手机号为空
		// 错误的收款账号
		{
			map[string]interface{}{
				"payChannelK": "3",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款人手机号不能为空",
			},
		},
		// 错误的收款账号
		{
			map[string]interface{}{
				"payChannelK": "3",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"phone":            "abcd",
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "收款人手机号不是一个有效的中国内地手机号码",
			},
		},
		// 请求金额需要大于1元
		{
			map[string]interface{}{
				"payChannelK": "3",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"identityCard":     "420222199212041057",
						"phone":            "13517210601",
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "0.01",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "50000",
				Message: "业务参数校验错误",
				SubCode: "api.common.error",
				SubMsg:  "付款金额金额不能少于1元",
			},
		},
		// 成功
		{
			map[string]interface{}{
				"payChannelK": "3",
				"orderData": [1]interface{}{
					map[string]interface{}{
						"identityCard":     "420222199212041057",
						"phone":            "13517210601",
						"orderSN":          uuid.New().String(),
						"payeeAccount":     "詹光",
						"requestPayAmount": "1",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
		// 多笔订单成功
		{
			map[string]interface{}{
				"payChannelK": "3",
				"orderData": [2]interface{}{
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"phone":            "13517210601",
						"payeeAccount":     "詹光",
						"requestPayAmount": "1",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
					map[string]interface{}{
						"orderSN":          uuid.New().String(),
						"phone":            "13517210601",
						"payeeAccount":     "詹光",
						"requestPayAmount": "1",
						"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					},
				},
			},
			cass.Response{
				HTTP:    cass.ResponseHTTP{StatusCode: 200},
				Code:    "10000",
				Message: "请求成功",
				SubCode: "",
				SubMsg:  "",
			},
		},
	}

	for _, tt := range tests {
		request := f.NewRequest(method.M.PayWeChatRemit)
		request.SetBizParams(tt.params)
		response := f.Send(request).(*cass.Response)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, tt.expectResp.HTTP.StatusCode, response.HTTP.StatusCode)
		assert.Equal(t, tt.expectResp.Code, response.Code)
		assert.Equal(t, tt.expectResp.SubCode, response.SubCode)
		assert.Equal(t, tt.expectResp.Message, response.Message)
		assert.Equal(t, tt.expectResp.SubMsg, response.SubMsg)
		t.Log(response.String())
	}
}

func TestMaYiBankPayNotAllowPayeeChannelType(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(method.M.PayOneBankRemit)
	request.SetBizParams(map[string]interface{}{
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
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTP.StatusCode)
	assert.Equal(t, "已选的属性 收款通道 非法", response.SubMsg)
	t.Log(response.String())
}

func TestOneWeChatPay(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(method.M.PayOneWeChatRemit)
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
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTP.StatusCode)
	t.Log(response.String())
}

func TestPayToMaYi(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(method.M.PayBankRemit)
	request.SetBizParams(map[string]interface{}{
		"payChannelK":      "1",
		"payeeChannelType": "1",
		"orderData": [2]interface{}{
			map[string]interface{}{
				"orderSN":          uuid.New().String(),
				"receiptFANO":      "6214860121221101",
				"payeeAccount":     "李浩然",
				"requestPayAmount": "80000",
				"notifyUrl":        "http://www.baidu.com/a/b?a=b",
				"identityCard":     "370404198511015755",
			},
			map[string]interface{}{
				"orderSN":          uuid.New().String(),
				"receiptFANO":      "6214860121221101",
				"payeeAccount":     "李浩然",
				"requestPayAmount": "90000",
				"notifyUrl":        "http://www.baidu.com/a/b?a=b",
				"identityCard":     "370404198511015755",
			},
		},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTP.StatusCode)
	t.Log(response.String())
}

func TestGetUsersVerifyStatus(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(method.M.GetUserVerifyStatus)
	request.SetBizParams(map[string]interface{}{
		"payChannelK":      "1",
		"payeeChannelType": "1",
		"orderData": [1]interface{}{
			map[string]interface{}{
				"orderSN":          uuid.New().String(),
				"receiptFANO":      "6214850278508756",
				"payeeAccount":     "13517210601",
				"identityCard":     "420222199212041057",
				"requestPayAmount": "0.01",
				"notifyUrl":        "http://www.baidu.com/a/b?a=b",
				"data":             "",
			},
		},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTP.StatusCode)
	t.Log(response.String())
}
