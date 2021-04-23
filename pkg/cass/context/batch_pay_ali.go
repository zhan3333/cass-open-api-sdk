package context

import "go-skysharing-openapi/pkg/cass/method"

var BatchPayALi = method.Method{
	Method: "Vzhuo.AliRemit.Pay",
	Name:   "批量支付宝实时下单",
}

type BatchPayALiBiz struct {
	Biz
	PayChannelK string     `json:"payChannelK"`
	OrderData   []ALiOrder `json:"orderData"`
}

type BatchPayALiContent struct {
	Content
	RBUUID string `json:"rbUUID"`
}
