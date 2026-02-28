package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cruisebooking/backend/internal/domain"
)

// 支付状态常量。
const (
	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusFailed  = "failed"
)

// PaymentGateway 向提供商发起支付订单。
// 返回提供商的交易号和向用户展示的支付链接。
type PaymentGateway interface {
	CreatePay(orderID int64, amountCents int64) (tradeNo string, payURL string, err error)
}

// PaymentService 处理支付创建。
type PaymentService struct {
	gw      PaymentGateway
	payRepo domain.PaymentRepository
}

// NewPaymentService 创建一个 PaymentService。
func NewPaymentService(gw PaymentGateway, payRepo domain.PaymentRepository) *PaymentService {
	return &PaymentService{gw: gw, payRepo: payRepo}
}

// Create 向网关发起支付订单，持久化待处理的支付记录，
// 并返回重定向用户的支付链接。
func (s *PaymentService) Create(ctx context.Context, orderID, amountCents int64, provider string) (string, error) {
	tradeNo, payURL, err := s.gw.CreatePay(orderID, amountCents)
	if err != nil {
		return "", fmt.Errorf("gateway.CreatePay: %w", err)
	}
	p := &domain.Payment{
		OrderID:     orderID,
		Provider:    provider,
		TradeNo:     tradeNo,
		AmountCents: amountCents,
		Status:      PaymentStatusPending,
	}
	if err := s.payRepo.Create(ctx, p); err != nil {
		return "", fmt.Errorf("persist payment: %w", err)
	}
	return payURL, nil
}

// ─── 回调验证 ────────────────────────────────────────────────────

// PaymentVerifier 验证来自支付提供商的已签名回调。
type PaymentVerifier interface {
	// Verify 根据提供商提供的签名验证原始回调主体。
	Verify(body []byte, signature string) error
	// ExtractTradeNo 从回调主体中解析提供商的交易号。
	ExtractTradeNo(body []byte) (string, error)
}

// HMACVerifier 使用 HMAC-SHA256 实现 PaymentVerifier。
//
// 生产环境提示：微信支付 v3 使用 RSA-SHA256 和平台公钥证书；
// 支付宝使用 RSA2 (SHA256withRSA)。在上线前请替换为相应 SDK 的实现。
type HMACVerifier struct{ secret string }

// NewHMACVerifier 使用给定的共享密钥创建一个 HMACVerifier。
func NewHMACVerifier(secret string) *HMACVerifier { return &HMACVerifier{secret: secret} }

// Verify 计算 HMAC-SHA256(body, secret) 并与提供的签名进行比较。
func (v *HMACVerifier) Verify(body []byte, signature string) error {
	if signature == "" {
		return errors.New("missing signature")
	}
	mac := hmac.New(sha256.New, []byte(v.secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	// 使用 hmac.Equal 防止时序攻击。
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return errors.New("signature mismatch")
	}
	return nil
}

// ExtractTradeNo 解析 JSON 主体并返回 "trade_no" 字段。
func (v *HMACVerifier) ExtractTradeNo(body []byte) (string, error) {
	var payload struct {
		TradeNo string `json:"trade_no"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("parse callback body: %w", err)
	}
	if payload.TradeNo == "" {
		return "", errors.New("trade_no is required in callback payload")
	}
	return payload.TradeNo, nil
}

// ─── 回调服务 ─────────────────────────────────────────────────────────

// BookingStatusUpdater 是回调服务所需的最小预订功能。
type BookingStatusUpdater interface {
	UpdateStatus(ctx context.Context, id int64, status string) error
}

// PaymentCallbackServiceImpl 处理异步的支付提供商回调。
type PaymentCallbackServiceImpl struct {
	payRepo     domain.PaymentRepository
	bookingRepo BookingStatusUpdater
	verifiers   map[string]PaymentVerifier
}

// NewPaymentCallbackService 创建一个 PaymentCallbackServiceImpl。
func NewPaymentCallbackService(
	payRepo domain.PaymentRepository,
	bookingRepo BookingStatusUpdater,
	verifiers map[string]PaymentVerifier,
) *PaymentCallbackServiceImpl {
	return &PaymentCallbackServiceImpl{
		payRepo:     payRepo,
		bookingRepo: bookingRepo,
		verifiers:   verifiers,
	}
}

// HandleCallback 处理支付提供商回调。可以针对同一个 trade_no 多次调用（幂等）。
//
// 流程：验证签名 → 提取 trade_no → 幂等性检查 →
//
//	更新支付状态 → 更新预订状态。
func (s *PaymentCallbackServiceImpl) HandleCallback(ctx context.Context, provider string, body []byte, signature string) error {
	v, ok := s.verifiers[provider]
	if !ok {
		return fmt.Errorf("unknown payment provider: %q", provider)
	}

	// 步骤 1：验证签名以拒绝伪造。
	if err := v.Verify(body, signature); err != nil {
		return fmt.Errorf("callback verification failed: %w", err)
	}

	// 步骤 2：提取提供商交易号。
	tradeNo, err := v.ExtractTradeNo(body)
	if err != nil {
		return err
	}

	// 步骤 3：幂等性 — 如果已支付，则返回成功且无副作用。
	payment, err := s.payRepo.FindByTradeNo(ctx, tradeNo)
	if err != nil {
		return fmt.Errorf("find payment by trade_no %q: %w", tradeNo, err)
	}
	if payment.Status == PaymentStatusPaid {
		return nil
	}

	// 步骤 4：将支付标记为已支付。
	if err := s.payRepo.UpdateStatus(ctx, payment.ID, PaymentStatusPaid); err != nil {
		return fmt.Errorf("update payment status: %w", err)
	}

	// 步骤 5：确认关联的预订。
	if err := s.bookingRepo.UpdateStatus(ctx, payment.OrderID, "paid"); err != nil {
		return fmt.Errorf("update booking status: %w", err)
	}

	return nil
}
