package context

import "go-skysharing-openapi/pkg/cass/method"

var QueryBatch = method.Method{
	Method: "Vzhuo.OneRemitStatus.Get",
	Name:   "获取单个批次状态",
}

type QueryBatchBiz struct {
	Biz
	RBUUID string `json:"rbUUID"`
}

type QueryBatchContent struct {
	Content
	RBUUID             string       `json:"rbUUID"`
	SBSNCN             string       `json:"SBSNCN"`
	TotalExpectAmount  string       `json:"totalExpectAmount"`
	TotalRealPayAmount string       `json:"totalRealPayAmount"`
	TotalServiceCharge string       `json:"totalServiceCharge"`
	DiscountAmount     string       `json:"discountAmount"`
	Status             int          `json:"status"`
	SubStatus          string       `json:"subStatus"`
	ResponseMsg        string       `json:"responseMsg"`
	OrderData          []BatchOrder `json:"orderData"`
}

type BatchOrder struct {
	OrderSN          string `json:"orderSN"`
	Phone            string `json:"phone"`
	OrderUUID        string `json:"orderUUID"`
	OrderStatus      int    `json:"orderStatus"`
	OrderFailStatus  string `json:"orderFailStatus"`
	RequestPayAmount string `json:"requestPayAmount"`
	ActualPayAmount  string `json:"actualPayAmount"`
	ReachAt          string `json:"reachAt"`
	OrderResponseMsg string `json:"orderResponseMsg"`
}
