package go_openapi_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go-skysharing-openapi/pkg/cass"
	"go-skysharing-openapi/pkg/cass/context"
	"sync"
	"testing"
	"time"
)

func TestOneAliPay(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(context.PayOneALi)
	request.SetBizParams(context.PayOneALiBiz{
		PayChannelK: context.PayChannelAliPay,
		OrderData: []context.ALiOrder{{
			OrderSN:          uuid.New().String(),
			ReceiptFANO:      "13517210601",
			PayeeAccount:     "詹光",
			RequestPayAmount: "0.01",
			NotifyUrl:        "https://www.baidu.com",
		}},
	})
	response := f.Send(request).(*cass.Response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.IsSuccess())
	c := &context.PayOneALiContent{}
	err = response.Content(c)
	assert.Nil(t, err)
	assert.NotEmpty(t, c.RBUUID)
}

// 测试批量
func TestBatchAliPay(t *testing.T) {
	f, err := cass.NewClient(factoryConf)
	assert.Nil(t, err)
	request := f.NewRequest(context.BatchPayALi)
	request.SetBizParams(context.BatchPayALiBiz{
		PayChannelK: context.PayChannelAliPay,
		OrderData: []context.ALiOrder{{
			OrderSN:          uuid.New().String(),
			ReceiptFANO:      "13517210601",
			PayeeAccount:     "詹光",
			RequestPayAmount: "0.01",
			NotifyUrl:        "https://www.baidu.com",
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

// 并发测试
func TestConcurrentOneALiPay(t *testing.T) {
	// 例如 seconds=60, TPS=10 时, 指的是每秒发起 10 次请求, 持续60秒, 一共发起 600 次请求
	var (
		// 指程序会运行多少秒
		seconds = 1
		// 指每秒会发起的请求数
		TPS              = 1
		wg               sync.WaitGroup
		successResponses []*cass.Response
		failedResponses  []*cass.Response
	)
	payOne := func(wg *sync.WaitGroup) {
		defer wg.Done()
		f, err := cass.NewClient(factoryConf)
		assert.Nil(t, err)
		request := f.NewRequest(context.PayOneALi)
		request.SetBizParams(context.PayOneALiBiz{
			PayChannelK: context.PayChannelAliPay,
			OrderData: []context.ALiOrder{{
				OrderSN:          uuid.New().String(),
				ReceiptFANO:      "18807159036",
				PayeeAccount:     "胡媛媛",
				RequestPayAmount: "0.01",
				NotifyUrl:        "https://www.baidu.com",
			}},
		})
		resp := f.Send(request).(*cass.Response)
		if resp == nil || !resp.IsSuccess() {
			failedResponses = append(failedResponses, resp)
		} else {
			successResponses = append(successResponses, resp)
		}
	}

	for i := 0; i < seconds; i++ {
		fmt.Printf("第 %d 波发起开始", i+1)
		for j := 0; j < TPS; j++ {
			wg.Add(1)
			go payOne(&wg)
			fmt.Printf("发起第 %d 波 第 %d 笔订单", i+1, j+1)
		}
		time.Sleep(time.Second * 1)
	}
	wg.Wait()

	fmt.Println("订单发起任务处理完毕")
	fmt.Println()
	fmt.Printf("共成功 %d 次请求\n", len(successResponses))
	fmt.Printf("共失败 %d 次请求\n", len(failedResponses))
	fmt.Println("失败订单如下")
	for i := range failedResponses {
		resp := failedResponses[i]
		fmt.Printf("- %d %s", resp.HTTP.StatusCode, resp)
	}
	fmt.Println()
	fmt.Println("程序执行完成退出")
}
