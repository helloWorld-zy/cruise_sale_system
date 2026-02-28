package service

import "fmt"

// WechatGateway 是微信支付网关的存根实现。
// 生产环境：请替换为真实的微信支付 v3 API（JSAPI + Native）集成。
type WechatGateway struct{}

// CreatePay 返回存根的交易号和支付链接。
func (w WechatGateway) CreatePay(orderID int64, _ int64) (string, string, error) {
	// TODO: 调用微信支付 v3 /v3/pay/transactions/jsapi 或 /native
	return fmt.Sprintf("wx_%d", orderID), "wechat://pay", nil
}

// AlipayGateway 是支付宝网关的存根实现。
// 生产环境：请替换为真实的支付宝 SDK (alipay.trade.page.pay) 集成。
type AlipayGateway struct{}

// CreatePay 返回存根的交易号和支付链接。
func (a AlipayGateway) CreatePay(orderID int64, _ int64) (string, string, error) {
	// TODO: 调用支付宝 alipay.trade.page.pay 或 alipay.trade.precreate
	return fmt.Sprintf("ali_%d", orderID), "alipay://pay", nil
}
