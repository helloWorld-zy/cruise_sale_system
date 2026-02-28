package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cruisebooking/backend/internal/domain"
)

// 退款状态常量。
const (
	RefundStatusPending  = "pending"
	RefundStatusApproved = "approved"
)

// RefundService 强制执行退款业务规则：
//   - amountCents 必须为正数
//   - 支付记录必须存在且处于 "paid" 状态
//   - amountCents 不能超过 (originalAmount − totalAlreadyRefunded)
type RefundService struct {
	payRepo    domain.PaymentRepository
	refundRepo domain.RefundRepository
}

// NewRefundService 创建一个 RefundService。
func NewRefundService(payRepo domain.PaymentRepository, refundRepo domain.RefundRepository) *RefundService {
	return &RefundService{payRepo: payRepo, refundRepo: refundRepo}
}

// Create 验证并持久化退款请求。
// 如果违反任何业务规则，则返回描述性错误。
func (s *RefundService) Create(ctx context.Context, paymentID, amountCents int64, reason string) error {
	if amountCents <= 0 {
		return errors.New("refund amount must be positive")
	}

	// 获取原始支付记录并确认已支付。
	payment, err := s.payRepo.FindByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}
	if payment.Status != PaymentStatusPaid {
		return fmt.Errorf("payment %d is not in paid status (current: %s)", paymentID, payment.Status)
	}

	// 计算已退款总额，以强制执行累计上限。
	alreadyRefunded, err := s.refundRepo.SumByPaymentID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("sum existing refunds: %w", err)
	}
	if alreadyRefunded+amountCents > payment.AmountCents {
		return fmt.Errorf(
			"refund amount %d exceeds remaining refundable balance %d (original=%d, already_refunded=%d)",
			amountCents, payment.AmountCents-alreadyRefunded, payment.AmountCents, alreadyRefunded,
		)
	}

	return s.refundRepo.Create(ctx, &domain.Refund{
		PaymentID:   paymentID,
		AmountCents: amountCents,
		Reason:      reason,
		Status:      RefundStatusPending,
	})
}
