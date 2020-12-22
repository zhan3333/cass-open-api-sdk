package context

import (
	"go-skysharing-openapi/pkg/cass/method"
)

const PayChannelBank = "1"
const PayChannelAliPay = "2"
const PayChannelWeChat = "3"
const PayeeChannelBank = "1"
const PayeeChannelAliPay = "2"

var GetBalance = method.Method{
	Method: "Vzhuo.BcBalance.Get",
	Name:   "获取余额",
}

type GetBalanceBiz struct {
	Biz         `json:"-"`
	PayChannelK string `json:"payChannelK"`
}

type GetBalanceContent struct {
	Content
	Bank struct {
		LockedAmount string `json:"lockedAmt"`
		CanUseAmount string `json:"canUseAmt"`
		Balance      string `json:"childFAbalance"`
	} `json:"bank"`
	AliPay struct {
		LockedAmount string `json:"lockedAmt"`
		CanUseAmount string `json:"canUseAmt"`
		Balance      string `json:"childFAbalance"`
	} `json:"alipay"`
	WeChat struct {
		LockedAmount string `json:"lockedAmt"`
		CanUseAmount string `json:"canUseAmt"`
		Balance      string `json:"childFAbalance"`
	} `json:"wechat"`
}
