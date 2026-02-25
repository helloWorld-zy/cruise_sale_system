package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/middleware"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// BookingService 定义预订处理器依赖的业务能力。
type BookingService interface {
	Create(userID, voyageID, skuID int64, guests int) error
}

// BookingHandler 处理 C 端预订下单请求。
type BookingHandler struct{ svc BookingService }

// NewBookingHandler 创建预订处理器实例。
func NewBookingHandler(svc BookingService) *BookingHandler { return &BookingHandler{svc: svc} }

// CreateBookingRequest 表示创建预订请求体。
type CreateBookingRequest struct {
	VoyageID   int64 `json:"voyage_id" binding:"required,gt=0"`
	CabinSKUID int64 `json:"cabin_sku_id" binding:"required,gt=0"`
	Guests     int   `json:"guests" binding:"required,gt=0"`
}

// Create 校验请求并创建预订。
func (h *BookingHandler) Create(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	if h.svc == nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "booking service unavailable")
		return
	}

	userValue, exists := c.Get(middleware.ContextKeyUserID) // M-01: C端订单使用 ContextKeyUserID
	if !exists {
		response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "not authenticated")
		return
	}

	userID, err := strconv.ParseInt(fmt.Sprint(userValue), 10, 64)
	if err != nil || userID <= 0 {
		response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "invalid user identity")
		return
	}

	if err := h.svc.Create(userID, req.VoyageID, req.CabinSKUID, req.Guests); err != nil {
		response.Error(c, http.StatusConflict, errcode.ErrConflict, err.Error())
		return
	}

	response.Success(c, gin.H{"status": "created"})
}
