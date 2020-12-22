package go_openapi_test

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go-skysharing-openapi/pkg/cass"
	"go-skysharing-openapi/pkg/cass/context"
	"go-skysharing-openapi/pkg/cass/method"
	"os"
	"strconv"
	"strings"
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
		PayChannelK: context.PayChannelBank,
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, response.Error())
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTP.StatusCode)
	assert.Equal(t, "10000", response.Code)
	assert.Equal(t, "请求成功", response.Message)
	assert.Equal(t, "", response.SubCode)
	assert.Equal(t, "", response.SubMsg)
	c := &context.GetBalanceContent{}
	err = response.Content(c)
	assert.Nil(t, err)
	assert.NotEmpty(t, c.Bank.Balance)
	assert.NotEmpty(t, c.Bank.CanUseAmount)
	assert.NotEmpty(t, c.Bank.LockedAmount)
}

// 测试单笔银行卡付款
func TestOneBankPay(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(context.PayOneBank)
	//request.SetBizParams(context.PayOneBankBiz{
	//	PayChannelK:      context.PayChannelBank,
	//	PayeeChannelType: context.PayeeChannelBank,
	//	OrderData: []context.PayOneBankOrder{{
	//		OrderSN:          uuid.New().String(),
	//		ReceiptFANO:      "6214850278508756",
	//		PayeeAccount:     "詹光",
	//		RequestPayAmount: "0.01",
	//		NotifyUrl:        "http://localhost:8080",
	//	}},
	//})
	request.SetBizParams(context.PayOneBankBiz{
		PayChannelK:      context.PayChannelBank,
		PayeeChannelType: context.PayeeChannelBank,
		OrderData: []context.PayOneBankOrder{{
			OrderSN:          uuid.New().String(),
			ReceiptFANO:      "6214850278508756",
			PayeeAccount:     "詹光",
			RequestPayAmount: "0.01",
			NotifyUrl:        "http://localhost:8080",
		}},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTP.StatusCode)
	t.Log(response.String())
	c := &context.PayOneBankContent{}
	err = response.Content(c)
	assert.Nil(t, err)
	t.Logf("%+v", c)
	assert.NotEmpty(t, c.RBUUID)
}

func TestOneAliPay(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(context.PayOneBank)
	request.SetBizParams(context.PayOneBankBiz{
		PayChannelK:      context.PayChannelBank,
		PayeeChannelType: context.PayeeChannelAliPay,
		OrderData: []context.PayOneBankOrder{{
			OrderSN:          uuid.New().String(),
			ReceiptFANO:      "13517210601",
			PayeeAccount:     "詹光",
			RequestPayAmount: "0.01",
			NotifyUrl:        "http://localhost:8080",
		}},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.IsSuccess())
	c := &context.PayOneBankContent{}
	err = response.Content(c)
	assert.Nil(t, err)
	assert.NotEmpty(t, c.RBUUID)
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
		request := f.NewRequest(context.PayOneBank)
		request.SetBizParams(tt.params)
		response := f.Send(request).(*cass.Response)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, tt.expectResp.HTTP.StatusCode, response.HTTP.StatusCode)
		assert.Equal(t, tt.expectResp.Code, response.Code, tt)
		assert.Equal(t, tt.expectResp.SubCode, response.SubCode)
		assert.Equal(t, tt.expectResp.Message, response.Message)
		assert.Equal(t, tt.expectResp.SubMsg, response.SubMsg)
		t.Log(response.String())
	}
}

// 测试批量提交网商银行支付
// 需要商户绑定通道为网商银行
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

	for k, tt := range tests {
		request := f.NewRequest(method.M.PayBankRemit)
		request.SetBizParams(tt.params)
		response := f.Send(request).(*cass.Response)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, tt.expectResp.HTTP.StatusCode, response.HTTP.StatusCode, map[string]interface{}{"tt": tt, "response": response.String(), "k": k})
		assert.Equal(t, tt.expectResp.Code, response.Code, map[string]interface{}{"tt": tt, "response": response.String(), "k": k})
		assert.Equal(t, tt.expectResp.SubCode, response.SubCode, map[string]interface{}{"tt": tt, "response": response.String(), "k": k})
		assert.Equal(t, tt.expectResp.Message, response.Message, map[string]interface{}{"tt": tt, "response": response.String(), "k": k})
		assert.Equal(t, tt.expectResp.SubMsg, response.SubMsg, map[string]interface{}{"tt": tt, "response": response.String(), "k": k})
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

// 测试网商银行通道不允许的收款类型
func TestMaYiBankPayNotAllowPayeeChannelType(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(context.PayOneBank)
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
	assert.True(t, response.IsHTTPSuccess())
	assert.False(t, response.IsBusinessSuccess())
	assert.Equal(t, "已选的属性收款通道非法", response.SubMsg, response.String())
}

// 测试单笔微信支付
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

// 查询批次
func TestQueryBatch(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(context.PayOneBank)
	request.SetBizParams(context.PayOneBankBiz{
		PayChannelK:      context.PayChannelBank,
		PayeeChannelType: context.PayeeChannelBank,
		OrderData: []context.PayOneBankOrder{{
			OrderSN:          uuid.New().String(),
			ReceiptFANO:      "6214850278508756",
			PayeeAccount:     "詹光",
			RequestPayAmount: "0.01",
			NotifyUrl:        "http://localhost:8080",
		}},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.True(t, response.IsSuccess())
	c := &context.PayOneBankContent{}
	assert.Nil(t, response.Content(c))

	query := f.NewRequest(context.QueryBatch)
	query.SetBizParams(context.QueryBatchBiz{
		RBUUID: c.RBUUID,
	})
	queryResp := f.Send(query).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, queryResp)
	batchC := &context.QueryBatchContent{}
	assert.Nil(t, queryResp.Content(batchC))
	assert.NotEmpty(t, batchC.RBUUID)
	//assert.NotEmpty(t, batchC.SBSNCN)
	assert.NotEmpty(t, batchC.TotalExpectAmount)
	assert.NotEmpty(t, batchC.TotalRealPayAmount)
	assert.NotEmpty(t, batchC.TotalServiceCharge)
	assert.NotEmpty(t, batchC.DiscountAmount)
	//assert.NotEmpty(t, batchC.Status)
	//assert.NotEmpty(t, batchC.SubStatus)
	//assert.NotEmpty(t, batchC.ResponseMsg)
	assert.NotEmpty(t, batchC.OrderData)
	assert.Equal(t, 1, len(batchC.OrderData))
	for _, order := range batchC.OrderData {
		assert.NotEmpty(t, order.OrderSN)
		//assert.NotEmpty(t, order.Phone)
		assert.NotEmpty(t, order.OrderUUID)
		//assert.NotEmpty(t, order.OrderStatus)
		//assert.NotEmpty(t, order.OrderFailStatus)
		assert.NotEmpty(t, order.RequestPayAmount)
		assert.NotEmpty(t, order.ActualPayAmount)
		//assert.NotEmpty(t, order.ReachAt)
		//assert.NotEmpty(t, order.OrderResponseMsg)
	}
}

// 通过 orderUUID 查询订单
func TestQueryOrderByUUID(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(context.PayOneBank)
	request.SetBizParams(context.PayOneBankBiz{
		PayChannelK:      context.PayChannelBank,
		PayeeChannelType: context.PayeeChannelBank,
		OrderData: []context.PayOneBankOrder{{
			OrderSN:          uuid.New().String(),
			ReceiptFANO:      "6214850278508756",
			PayeeAccount:     "詹光",
			RequestPayAmount: "0.01",
			NotifyUrl:        "http://localhost:8080",
		}},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.True(t, response.IsSuccess())
	c := &context.PayOneBankContent{}
	assert.Nil(t, response.Content(c))

	query := f.NewRequest(context.QueryBatch)
	query.SetBizParams(context.QueryBatchBiz{
		RBUUID: c.RBUUID,
	})
	queryResp := f.Send(query).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, queryResp)
	batchC := &context.QueryBatchContent{}
	assert.Nil(t, queryResp.Content(batchC))
	assert.NotEmpty(t, batchC.RBUUID)
	assert.NotEmpty(t, batchC.TotalExpectAmount)
	assert.NotEmpty(t, batchC.TotalRealPayAmount)
	assert.NotEmpty(t, batchC.TotalServiceCharge)
	assert.NotEmpty(t, batchC.DiscountAmount)
	assert.NotEmpty(t, batchC.OrderData)
	assert.Equal(t, 1, len(batchC.OrderData))
	for _, order := range batchC.OrderData {
		assert.NotEmpty(t, order.OrderSN)
		assert.NotEmpty(t, order.OrderUUID)
		assert.NotEmpty(t, order.RequestPayAmount)
		assert.NotEmpty(t, order.ActualPayAmount)
	}

	var orderUUID = batchC.OrderData[0].OrderUUID
	queryOrderReq := f.NewRequest(context.QueryOrderByUUID)
	queryOrderReq.SetBizParams(context.QueryOrderByUUIDBiz{
		OrderUUID: orderUUID,
	})
	queryOrderResp := f.Send(queryOrderReq)
	assert.Nil(t, queryOrderResp.Error())
	t.Log(queryOrderResp.String())
	queryOrderC := &context.QueryOrderByUUIDContent{}
	assert.Nil(t, queryOrderResp.Content(queryOrderC))
	assert.NotEmpty(t, queryOrderC.RBUUID)
	assert.NotEmpty(t, queryOrderC.OrderUUID)
	assert.NotEmpty(t, queryOrderC.OrderSN)
}

// 通过 orderSN 查询订单
func TestQueryOrderBySN(t *testing.T) {
	var orderSN = strings.ReplaceAll(uuid.New().String(), "-", "")
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(context.PayOneBank)
	request.SetBizParams(context.PayOneBankBiz{
		PayChannelK:      context.PayChannelBank,
		PayeeChannelType: context.PayeeChannelBank,
		OrderData: []context.PayOneBankOrder{{
			OrderSN:          orderSN,
			ReceiptFANO:      "6214850278508756",
			PayeeAccount:     "詹光",
			RequestPayAmount: "0.01",
			NotifyUrl:        "http://localhost:8080",
		}},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.True(t, response.IsSuccess())
	c := &context.PayOneBankContent{}
	assert.Nil(t, response.Content(c))

	queryOrderReq := f.NewRequest(context.QueryOrderBySN)
	queryOrderReq.SetBizParams(context.QueryOrderBySNBiz{
		OrderSN: orderSN,
	})
	queryOrderResp := f.Send(queryOrderReq)
	assert.Nil(t, queryOrderResp.Error())
	t.Log(queryOrderResp.String())
	queryOrderC := &context.QueryOrderBySNContent{}
	assert.Nil(t, queryOrderResp.Content(queryOrderC))
	assert.NotEmpty(t, queryOrderC.RBUUID)
	assert.NotEmpty(t, queryOrderC.OrderUUID)
	assert.NotEmpty(t, queryOrderC.OrderSN)
}
