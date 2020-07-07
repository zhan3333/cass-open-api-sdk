package cass_test

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/srlemon/gen-id"
	"github.com/stretchr/testify/assert"
	"go-skysharing-openapi/pkg/cass"
	"math/rand"
	"testing"
)

func BenchmarkPay(b *testing.B) {
	f, err := cass.NewFactory(factoryConf)
	assert.Nil(b, err)
	assert.NotNil(b, f)
	b.RunParallel(func(pb *testing.PB) {
		// 每个 goroutine 有属于自己的 bytes.Buffer.
		var buf bytes.Buffer
		var err error
		var response *cass.Response
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			buf.Reset()
			useMethods := []cass.Method{
				cass.M.PayOneBankRemit,
				cass.M.PayBankRemit,
				cass.M.PayOneWeChatRemit,
				cass.M.PayWeChatRemit,
			}
			method := useMethods[rand.Intn(4)]
			req := f.NewRequest(method)
			switch method.Num {
			case cass.M.PayOneBankRemit.Num:
				req.SetBizParams(cass.PayOneBankRemitBiz{
					PayChannelK:      "1",
					PayeeChannelType: "1",
					OrderData: []cass.BankOrder{
						{
							OrderSN:          uuid.New().String(),
							ReceiptFANO:      genid.NewGeneratorData().BankID,
							PayeeAccount:     genid.NewGeneratorData().Name,
							RequestPayAmount: "0.01",
							NotifyUrl:        "http://www.baidu.com/",
						},
					},
				})
				break
			case cass.M.PayOneWeChatRemit.Num:
				req.SetBizParams(cass.PayOneWeChatRemitBiz{
					PayChannelK: "3",
					OrderData: []cass.WeChatOrder{
						{
							OrderSN:          uuid.New().String(),
							Phone:            genid.NewGeneratorData().PhoneNum,
							PayeeAccount:     "詹光",
							RequestPayAmount: "1",
							NotifyUrl:        "http://www.baidu.com/",
						},
					},
				})
				break
			case cass.M.PayBankRemit.Num:
				req.SetBizParams(cass.PayBankRemitBiz{
					PayChannelK:      "1",
					PayeeChannelType: "1",
					OrderData: []cass.BankOrder{
						{
							OrderSN:          uuid.New().String(),
							ReceiptFANO:      genid.NewGeneratorData().BankID,
							PayeeAccount:     genid.NewGeneratorData().Name,
							RequestPayAmount: "0.01",
							NotifyUrl:        "http://www.baidu.com/",
						},
					},
				})
				break
			case cass.M.PayWeChatRemit.Num:
				req.SetBizParams(cass.PayOneWeChatRemitBiz{
					PayChannelK: "3",
					OrderData: []cass.WeChatOrder{
						{
							OrderSN:          uuid.New().String(),
							Phone:            genid.NewGeneratorData().PhoneNum,
							PayeeAccount:     "詹光",
							RequestPayAmount: "1",
							NotifyUrl:        "http://www.baidu.com/",
						},
					},
				})
				break
			default:
				b.Errorf("错误的method: %+v", method)
				b.FailNow()
			}
			response, err = req.Send()
			assert.Nil(b, err)
			//b.Logf("%s: %s", method.Name, response.String())
			assert.NotNil(b, response)
			assert.Equal(b, 200, response.StatusCode)
			assert.Equal(b, "10000", response.Code)
		}
	})
}

func BenchmarkPayBank(b *testing.B) {
	f, err := cass.NewFactory(factoryConf)
	assert.Nil(b, err)
	assert.NotNil(b, f)
	b.RunParallel(func(pb *testing.PB) {
		// 每个 goroutine 有属于自己的 bytes.Buffer.
		var buf bytes.Buffer
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
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
