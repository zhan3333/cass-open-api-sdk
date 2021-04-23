package cass

type PayBankRemitBiz struct {
	PayChannelK      string      `json:"payChannelK" binding:"required,max=1" comment:" 付款通道：1-银行卡"`
	PayeeChannelType string      `json:"payeeChannelType" binding:"max=1" comment:"收款通道： 网商银行为必填，1-银行卡，2-支付宝； 其他通道，收款与付款通道一致，不需要接 口传入数据。"`
	ContractID       string      `json:"contractID,omitempty" binding:"" comment:"合同 ID"`
	OrderData        []BankOrder `json:"orderData" binding:"required" comment:" 订单数据，二维数组"`
}

type BankOrder struct {
	OrderSN          string `json:"orderSN" binding:"required,max=64" comment:"商户订单号，只能是英文字母，数字，中文 以及连接符-" example:"测试商户-20190712-0026"`
	ReceiptFANO      string `json:"receiptFANO" binding:"required,max=23" comment:"收款人金融账号： 网商银行支持，1-银行账号 2-支付宝账号，邮箱，手机号 其他银行仅支持，1-银行账号" example:"银行卡：6214832719548955 支付宝：18800080008"`
	PayeeAccount     string `json:"payeeAccount" binding:"required,max=64" comment:"  收款人户名（真实姓名）" example:"张三"`
	ReceiptBankName  string `json:"receiptBankName,omitempty" binding:"max=64" comment:"收款方开户行名称（总行即可）" example:"招商银行"`
	ReceiptBankAddr  string `json:"receiptBankAddr,omitempty" binding:"max=100" comment:"  收款方开户行地（省会或者直辖市即可）" example:"武汉"`
	CRCHGNO          string `json:"CRCHGNO,omitempty" binding:"max=12" comment:"收方联行号" example:" 411234567222"`
	RequestPayAmount string `json:"requestPayAmount" binding:"required,max=26" comment:"预期付款金额" example:"14.00"`
	IdentityCard     string `json:"identityCard,omitempty" binding:"max=20" comment:"收款人身份证号" example:" 321123456789098765"`
	NotifyUrl        string `json:"notifyUrl" binding:"required,max=255" comment:"网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单”" example:" └notifyUrl String -- 255 网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单” http://xxx.xxx.cn/xx/asynNotify.h tm"`
}

type PayOneWeChatRemitBiz struct {
	PayChannelK string        `json:"payChannelK" binding:"required,max=1" comment:" 付款通道：3-微信"`
	ContractID  string        `json:"contractID,omitempty" binding:"" comment:"合同 ID"`
	OrderData   []WeChatOrder `json:"orderData" binding:"required" comment:" 订单数据，二维数组"`
}

type PayWeChatRemitBiz struct {
	PayChannelK string        `json:"payChannelK" binding:"required,max=1" comment:" 付款通道：3-微信"`
	ContractID  string        `json:"contractID,omitempty" binding:"" comment:"合同 ID"`
	OrderData   []WeChatOrder `json:"orderData" binding:"required" comment:" 订单数据，二维数组"`
}

type WeChatOrder struct {
	OrderSN          string `json:"orderSN" binding:"required,max=64" comment:"商户订单号，只能是英文字母，数字，中文 以及连接符-" example:"测试商户-20190712-0026"`
	Phone            string `json:"phone" binding:"required,max=12" comment:"收款人手机号" example:"13517588888"`
	PayeeAccount     string `json:"payeeAccount" binding:"required,max=64" comment:"  收款人户名（真实姓名）" example:"张三"`
	RequestPayAmount string `json:"requestPayAmount" binding:"required,max=26" comment:"预期付款金额" example:"14.00"`
	IdentityCard     string `json:"identityCard,omitempty" binding:"max=20" comment:"收款人身份证号" example:" 321123456789098765"`
	NotifyUrl        string `json:"notifyUrl" binding:"required,max=255" comment:"网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单”" example:" └notifyUrl String -- 255 网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单” http://xxx.xxx.cn/xx/asynNotify.h tm"`
}

type GetUserVerifyStatusItem struct {
	IdentityCard string `json:"identityCard"`
	ReceiptFANO  string `json:"receiptFANO"`
	ReceiptType  int    `json:"receiptType"`
}

type GetUserVerifyStatusBiz struct {
	Items []GetUserVerifyStatusItem `json:"items"`
}

// 支付宝订单
type ALiOrder struct {
	OrderSN          string `json:"orderSN" binding:"required,max=64" comment:"商户订单号，只能是英文字母，数字，中文 以及连接符-" example:"测试商户-20190712-0026"`
	ReceiptFANO      string `json:"receiptFANO" binding:"required,max=23" comment:"收款人金融账号： 网商银行支持，1-银行账号 2-支付宝账号，邮箱，手机号 其他银行仅支持，1-银行账号" example:"银行卡：6214832719548955 支付宝：18800080008"`
	PayeeAccount     string `json:"payeeAccount" binding:"required,max=64" comment:"  收款人户名（真实姓名）" example:"张三"`
	RequestPayAmount string `json:"requestPayAmount" binding:"required,max=26" comment:"预期付款金额" example:"14.00"`
	IdentityCard     string `json:"identityCard,omitempty" binding:"max=20" comment:"收款人身份证号" example:" 321123456789098765"`
	NotifyUrl        string `json:"notifyUrl" binding:"required,max=255" comment:"网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单”" example:" └notifyUrl String -- 255 网商银行必填，服务器异步通知页面路径。 对应异步通知的“银行卡实时下单” http://xxx.xxx.cn/xx/asynNotify.h tm"`
	Tax              string `json:"tax" binding:"optiempty" comment:"个税"`
}

type PayAliBiz struct {
	PayChannelK string     `json:"payChannelK"`
	Orders      []ALiOrder `json:"orderData"`
}

type BatchPayAliBiz struct {
	PayChannelK string     `json:"payChannelK"`
	Orders      []ALiOrder `json:"orderData"`
}
