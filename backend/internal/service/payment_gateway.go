package service

type WechatGateway struct{}

func (w WechatGateway) CreatePay(orderID int64, amountCents int64) (string, error) {
	return "wechat://pay", nil
}

type AlipayGateway struct{}

func (a AlipayGateway) CreatePay(orderID int64, amountCents int64) (string, error) {
	return "alipay://pay", nil
}
