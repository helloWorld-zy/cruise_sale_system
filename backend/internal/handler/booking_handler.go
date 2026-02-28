package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/middleware"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// BookingService 定义预订处理器依赖的业务能力。
type BookingService interface {
	Create(ctx context.Context, userID, voyageID, skuID int64, guests int) (*domain.Booking, error)
}

// BookingAdminStore 定义管理后台订单查询与管理能力。
type BookingAdminStore interface {
	List(ctx context.Context, page, pageSize int) ([]domain.Booking, int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	Delete(ctx context.Context, id int64) error
}

// BookingHandler 处理 C 端预订下单请求。
type BookingHandler struct {
	svc        BookingService
	adminStore BookingAdminStore
}

// NewBookingHandler 创建预订处理器实例。
func NewBookingHandler(svc BookingService, adminStore ...BookingAdminStore) *BookingHandler {
	h := &BookingHandler{svc: svc}
	if len(adminStore) > 0 {
		h.adminStore = adminStore[0]
	}
	return h
}

// CreateBookingRequest 表示创建预订请求体。
type CreateBookingRequest struct {
	UserID     int64 `json:"user_id"`
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

	var userID int64
	if userValue, exists := c.Get(middleware.ContextKeyUserID); exists {
		parsedID, err := strconv.ParseInt(fmt.Sprint(userValue), 10, 64)
		if err != nil || parsedID <= 0 {
			response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "invalid user identity")
			return
		}
		userID = parsedID
	} else if req.UserID > 0 {
		// 管理后台允许显式指定下单用户。
		userID = req.UserID
	} else {
		response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "not authenticated")
		return
	}

	if userID <= 0 {
		response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "invalid user identity")
		return
	}

	booking, err := h.svc.Create(c.Request.Context(), userID, req.VoyageID, req.CabinSKUID, req.Guests)
	if err != nil {
		response.Error(c, http.StatusConflict, errcode.ErrConflict, err.Error())
		return
	}

	response.Success(c, gin.H{"id": booking.ID, "status": booking.Status, "total_cents": booking.TotalCents})
}

// AdminList 管理后台分页查询订单。
func (h *BookingHandler) AdminList(c *gin.Context) {
	if h.adminStore == nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "booking store unavailable")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	items, total, err := h.adminStore.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, gin.H{"list": items, "total": total})
}

// AdminGet 管理后台查询订单详情。
func (h *BookingHandler) AdminGet(c *gin.Context) {
	if h.adminStore == nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "booking store unavailable")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	b, err := h.adminStore.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "booking not found")
		return
	}
	response.Success(c, b)
}

// AdminUpdate 管理后台更新订单状态。
func (h *BookingHandler) AdminUpdate(c *gin.Context) {
	if h.adminStore == nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "booking store unavailable")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if err := h.adminStore.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, gin.H{"id": id, "status": req.Status})
}

// AdminDelete 管理后台删除订单。
func (h *BookingHandler) AdminDelete(c *gin.Context) {
	if h.adminStore == nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "booking store unavailable")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.adminStore.Delete(c.Request.Context(), id); err != nil {
		respondDeleteError(c, err, "booking")
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateStatus 测试用的空实现
func (h *BookingHandler) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	return nil
}
