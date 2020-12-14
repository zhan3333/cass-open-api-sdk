package context

import "go-skysharing-openapi/pkg/cass/method"

var QueryOrderByUUID = method.Method{
	Method: "Vzhuo.OneOrderStatus.Get",
	Name:   "获取单个订单状态",
}

var QueryOrderBySN = method.Method{
	Method: "Vzhuo.OneOrderByOuterOrderSN.Get",
	Name:   "获取单个订单状态",
}

type QueryOrderByUUIDBiz struct {
	Biz
	OrderUUID string `json:"orderUUID"`
}

type QueryOrderByUUIDContent struct {
	Content
	RBUUID      string `json:"rbUUID"`
	OrderUUID   string `json:"orderUUID"`
	OrderSN     string `json:"orderSN"`
	OrderStatus int    `json:"orderStatus"`
	RemitStatus int    `json:"remitStatus"`
	ReachAt     string `json:"reachAt"`
	ResponseMsg string `json:"responseMsg"`
}

type QueryOrderBySNBiz struct {
	Biz
	OrderSN string `json:"orderSN"`
}

type QueryOrderBySNContent struct {
	Content
	RBUUID      string `json:"rbUUID"`
	OrderUUID   string `json:"orderUUID"`
	OrderSN     string `json:"orderSN"`
	OrderStatus int    `json:"orderStatus"`
	RemitStatus int    `json:"remitStatus"`
	ReachAt     string `json:"reachAt"`
	ResponseMsg string `json:"responseMsg"`
}
