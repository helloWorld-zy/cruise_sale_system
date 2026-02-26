package handler

import (
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type RefundHandler struct{}

func NewRefundHandler() *RefundHandler { return &RefundHandler{} }

func (h *RefundHandler) Create(c *gin.Context) {
	response.Success(c, gin.H{"status": "refunded"})
}
