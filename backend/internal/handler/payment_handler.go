package handler

import (
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type PaymentCallbackService interface{ HandleCallback(payload []byte) error }

type PaymentHandler struct{ svc PaymentCallbackService }

func NewPaymentHandler(svc PaymentCallbackService) *PaymentHandler { return &PaymentHandler{svc: svc} }

func (h *PaymentHandler) Callback(c *gin.Context) {
	response.Success(c, gin.H{"status": "paid"})
}
