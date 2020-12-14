package context

import "go-skysharing-openapi/pkg/cass/method"

var PayOneBank = method.Method{
	Method: "Vzhuo.OneBankRemit.Pay",
	Name:   "单笔银行卡实时下单",
}

type PayOneBankBiz struct {
	Biz
	PayChannelK      string            `json:"payChannelK"`
	PayeeChannelType string            `json:"payeeChannelType"`
	OrderData        []PayOneBankOrder `json:"orderData"`
}

type PayOneBankOrder struct {
	OrderSN          string        `json:"orderSN"`
	ReceiptFANO      string        `json:"receiptFANO"`
	PayeeAccount     string        `json:"payeeAccount"`
	RequestPayAmount string        `json:"requestPayAmount"`
	NotifyUrl        string        `json:"notifyUrl"`
	Data             []interface{} `json:"data"`
}

type PayOneBankContent struct {
	Content
	RBUUID string `json:"rbUUID"`
}
