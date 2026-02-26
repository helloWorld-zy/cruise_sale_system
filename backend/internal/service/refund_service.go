package service

import "github.com/cruisebooking/backend/internal/domain"

type RefundRepo interface{ Create(v interface{}) error }

type RefundService struct{ repo RefundRepo }

func NewRefundService(repo RefundRepo) *RefundService { return &RefundService{repo: repo} }

func (s *RefundService) Create(paymentID int64, amountCents int64, reason string) error {
	refund := domain.Refund{
		PaymentID:   paymentID,
		AmountCents: amountCents,
		Reason:      reason,
	}
	return s.repo.Create(&refund)
}
