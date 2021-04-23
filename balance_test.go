package go_openapi_test

import (
	"github.com/stretchr/testify/assert"
	"go-skysharing-openapi/pkg/cass"
	"go-skysharing-openapi/pkg/cass/context"
	"testing"
)

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
	t.Logf("%+v", c)
}
