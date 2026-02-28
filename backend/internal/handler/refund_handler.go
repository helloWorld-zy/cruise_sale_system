package handler

import (
	"context"
	"net/http"

	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// RefundCreateService 是退款处理器所需的依赖服务。
type RefundCreateService interface {
	Create(ctx context.Context, paymentID, amountCents int64, reason string) error
}

// RefundHandler 处理退款相关的 HTTP 请求。
type RefundHandler struct{ svc RefundCreateService }

// NewRefundHandler 使用给定的服务创建 RefundHandler 实例。
func NewRefundHandler(svc RefundCreateService) *RefundHandler { return &RefundHandler{svc: svc} }

// CreateRefundRequest 表示创建退款的请求体。
type CreateRefundRequest struct {
	PaymentID   int64  `json:"payment_id"   binding:"required,gt=0"`
	AmountCents int64  `json:"amount_cents" binding:"required,gt=0"`
	Reason      string `json:"reason"       binding:"required"`
}

// Create 处理 POST /api/v1/refunds 请求。
// 服务层会强制校验：amountCents ≤ (originalAmount − alreadyRefunded)。
func (h *RefundHandler) Create(c *gin.Context) {
	var req CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	if err := h.svc.Create(c.Request.Context(), req.PaymentID, req.AmountCents, req.Reason); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, errcode.ErrConflict, err.Error())
		return
	}

	response.Success(c, gin.H{"status": "refunded"})
}
