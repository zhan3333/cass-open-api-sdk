package context

import "go-skysharing-openapi/pkg/cass/method"

var PayOneALi = method.Method{
	Method: "Vzhuo.OneAliRemit.Pay",
	Name:   "单笔支付宝实时下单",
}

type PayOneALiBiz struct {
	Biz
	PayChannelK string     `json:"payChannelK"`
	OrderData   []ALiOrder `json:"orderData"`
}

type PayOneALiContent struct {
	Content
	RBUUID string `json:"rbUUID"`
}

type ALiOrder struct {
	OrderSN          string `json:"orderSN" binding:"required,max=64" comment:"商户订单号，只能是英文字母，数字，中文 以及连接符-" example:"测试商户-20190712-0026"`
	ReceiptFANO      string `json:"receiptFANO" binding:"required,max=23" comment:"收款人金融账号： 网商银行支持，1-银行账号 2-支付宝账号，邮箱，手机号 其他银行仅支持，1-银行账号" example:"银行卡：6214832719548955 支付宝：18800080008"`
	PayeeAccount     string `json:"payeeAccount" binding:"required,max=64" comment:"  收款人户名（真实姓名）" example:"张三"`
	RequestPayAmount string `json:"requestPayAmount" binding:"required,max=26" comment:"预期付款金额" example:"14.00"`
	IdentityCard     string `json:"identityCard,omitempty" binding:"max=20" comment:"收款人身份证号" example:" 321123456789098765"`
	NotifyUrl        string `json:"notifyUrl" binding:"required,max=255" comment:"网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单”" example:" └notifyUrl String -- 255 网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单” http://xxx.xxx.cn/xx/asynNotify.h tm"`
	Tax              string `json:"tax" binding:"optiempty" comment:"个税"`
}
