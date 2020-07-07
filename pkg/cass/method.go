package cass

type Methods struct {
	GetBalance           Method
	GetOneRemitStatus    Method
	GetOneOrderStatus    Method
	PayBankRemit         Method
	PayOneBankRemit      Method
	PayWeChatRemit       Method
	PayOneWeChatRemit    Method
	GetChannelData       Method
	ChargeBank           Method
	ChargeWeChat         Method
	GetApplyChargeResult Method
	VerifyUser           Method
	GetUserVerifyStatus  Method
	GetContractList      Method
}

var M = Methods{
	GetBalance: Method{
		Num:    "3.1",
		Method: "Vzhuo.BcBalance.Get",
		Name:   "获取余额",
	},
	GetOneRemitStatus: Method{
		Num:    "3.2",
		Method: "Vzhuo.OneRemitStatus.Get",
		Name:   "获取单个批次状态",
	},
	GetOneOrderStatus: Method{
		Num:    "3.3",
		Method: "Vzhuo.OneOrderStatus.Get",
		Name:   "获取单个订单状态",
	},
	PayBankRemit: Method{
		Num:    "3.4",
		Method: "Vzhuo.BankRemit.Pay",
		Name:   "银行卡实时下单",
	},
	PayWeChatRemit: Method{
		Num:    "3.5",
		Method: "Vzhuo.WeChatRemit.Pay",
		Name:   "微信实时下单",
	},
	PayOneBankRemit: Method{
		Num:    "3.6",
		Method: "Vzhuo.OneBankRemit.Pay",
		Name:   "单笔银行卡实时下单",
	},
	PayOneWeChatRemit: Method{
		Num:    "3.7",
		Method: "Vzhuo.OneWeChatRemit.Pay",
		Name:   "单笔微信实时下单",
	},
	GetChannelData: Method{
		Num:    "3.8",
		Method: "Vzhuo.ChannelData.Get",
		Name:   "充值账号查询",
	},
	ChargeBank: Method{
		Num:    "3.9",
		Method: "Vzhuo.Bank.Charge",
		Name:   "银行卡充值申请提交",
	},
	ChargeWeChat: Method{
		Num:    "3.10",
		Method: "Vzhuo.WeChat.Charge",
		Name:   "微信充值申请提交",
	},
	GetApplyChargeResult: Method{
		Num:    "3.11",
		Method: "Vzhuo.ApplyChargeResult.Get",
		Name:   "获取充值结果",
	},
	VerifyUser: Method{
		Num:    "3.12",
		Method: "Vzhuo.User.Verify",
		Name:   "添加用户实名认证结果",
	},
	GetUserVerifyStatus: Method{
		Num:    "3.13",
		Method: "Vzhuo.UsersVerifyStatus.Get",
		Name:   "批量查询用户认证状态",
	},
	GetContractList: Method{
		Num:    "3.14",
		Method: "Vzhuo.ContractList.Get",
		Name:   "获取合同列表",
	},
}

type Method struct {
	Num    string
	Method string
	Name   string
}

func (m Methods) GetOptions() []string {
	return []string{
		m.GetBalance.Method,
		m.GetOneRemitStatus.Method,
		m.GetOneOrderStatus.Method,
		m.PayBankRemit.Method,
		m.PayOneBankRemit.Method,
		m.GetChannelData.Method,
		m.ChargeBank.Method,
		m.GetApplyChargeResult.Method,
		m.VerifyUser.Method,
		m.GetUserVerifyStatus.Method,
		m.GetContractList.Method,
	}
}

func (m Methods) GetOption(name string) Method {
	switch name {
	case m.GetBalance.Method:
		return m.GetBalance
	case m.GetOneRemitStatus.Method:
		return m.GetOneRemitStatus
	case m.GetOneOrderStatus.Method:
		return m.GetOneOrderStatus
	case m.PayBankRemit.Method:
		return m.PayBankRemit
	case m.PayOneBankRemit.Method:
		return m.PayOneBankRemit
	case m.GetChannelData.Method:
		return m.GetChannelData
	case m.ChargeBank.Method:
		return m.ChargeBank
	case m.GetApplyChargeResult.Method:
		return m.GetApplyChargeResult
	case m.VerifyUser.Method:
		return m.VerifyUser
	case m.GetUserVerifyStatus.Method:
		return m.GetUserVerifyStatus
	case m.GetContractList.Method:
		return m.GetContractList
	default:
		return m.GetBalance
	}
}
