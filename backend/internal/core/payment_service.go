package core

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) GetPaymentURL(orderNo string, amount float64) string {
	// Mock URL
	return "https://mock-payment-gateway.com/pay?order=" + orderNo
}
