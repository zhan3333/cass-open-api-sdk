package method

import (
	"fmt"
	"strings"
)

type Methods struct {
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
	GetOneRemitStatus: Method{
		Method: "Vzhuo.OneRemitStatus.Get",
		Name:   "获取单个批次状态",
	},
	GetOneOrderStatus: Method{
		Method: "Vzhuo.OneOrderStatus.Get",
		Name:   "获取单个订单状态",
	},
	PayBankRemit: Method{
		Method: "Vzhuo.BankRemit.Pay",
		Name:   "银行卡实时下单",
	},
	PayWeChatRemit: Method{
		Method: "Vzhuo.WeChatRemit.Pay",
		Name:   "微信实时下单",
	},
	PayOneBankRemit: Method{
		Method: "Vzhuo.OneBankRemit.Pay",
		Name:   "单笔银行卡实时下单",
	},
	PayOneWeChatRemit: Method{
		Method: "Vzhuo.OneWeChatRemit.Pay",
		Name:   "单笔微信实时下单",
	},
	GetChannelData: Method{
		Method: "Vzhuo.ChannelData.Get",
		Name:   "充值账号查询",
	},
	ChargeBank: Method{
		Method: "Vzhuo.Bank.Charge",
		Name:   "银行卡充值申请提交",
	},
	ChargeWeChat: Method{
		Method: "Vzhuo.WeChat.Charge",
		Name:   "微信充值申请提交",
	},
	GetApplyChargeResult: Method{
		Method: "Vzhuo.ApplyChargeResult.Get",
		Name:   "获取充值结果",
	},
	VerifyUser: Method{
		Method: "Vzhuo.User.Verify",
		Name:   "添加用户实名认证结果",
	},
	GetUserVerifyStatus: Method{
		Method: "Vzhuo.UsersVerifyStatus.Get",
		Name:   "批量查询用户认证状态",
	},
	GetContractList: Method{
		Method: "Vzhuo.ContractList.Get",
		Name:   "获取合同列表",
	},
}

type Method struct {
	Method string
	Name   string
}

func (m Method) GetResponseKey() string {
	s := strings.ReplaceAll(m.Method, ".", "")
	return fmt.Sprintf("%sResponse", s)
}
