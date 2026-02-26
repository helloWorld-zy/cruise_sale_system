package service

type PaymentGateway interface {
	CreatePay(orderID int64, amountCents int64) (string, error)
}

type PaymentService struct{ gw PaymentGateway }

func NewPaymentService(gw PaymentGateway) *PaymentService { return &PaymentService{gw: gw} }

func (s *PaymentService) Create(orderID int64, amountCents int64) (string, error) {
	return s.gw.CreatePay(orderID, amountCents)
}
