package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/middleware"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/repository"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// BookingService 定义预订处理器依赖的业务能力。
type BookingService interface {
	Create(ctx context.Context, userID, voyageID, skuID int64, guests int) (*domain.Booking, error)
}

// BookingAdminStore 定义管理后台订单查询与管理能力。
type BookingAdminStore interface {
	List(ctx context.Context, page, pageSize int) ([]domain.Booking, int64, error)
	ListWithFilter(ctx context.Context, filter repository.BookingFilter, page, pageSize int) ([]domain.Booking, int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	TransitionStatus(ctx context.Context, id int64, status string, operatorID int64, remark string) error
	Delete(ctx context.Context, id int64) error
}

// BookingHandler 处理 C 端预订下单请求。
type BookingHandler struct {
	svc           BookingService
	adminStore    BookingAdminStore
	exportService *service.OrderExportService
}

// NewBookingHandler 创建预订处理器实例。
func NewBookingHandler(svc BookingService, adminStore ...BookingAdminStore) *BookingHandler {
	h := &BookingHandler{svc: svc}
	if len(adminStore) > 0 {
		h.adminStore = adminStore[0]
	}
	return h
}

// SetExportService 注入订单导出服务。
func (h *BookingHandler) SetExportService(exportSvc *service.OrderExportService) {
	h.exportService = exportSvc
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

	filter := repository.BookingFilter{
		Status:    c.Query("status"),
		VoyageID:  0,
		BookingNo: c.Query("booking_no"),
		StartDate: nil,
		EndDate:   nil,
	}

	if voyageID, err := strconv.ParseInt(c.Query("voyage_id"), 10, 64); err == nil {
		filter.VoyageID = voyageID
	}
	if startDate := c.Query("start_date"); startDate != "" {
		filter.StartDate = &startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filter.EndDate = &endDate
	}

	items, total, err := h.adminStore.ListWithFilter(c.Request.Context(), filter, page, pageSize)
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
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	operatorID := parseOperatorID(c)
	if req.Remark == "" {
		req.Remark = fmt.Sprintf("admin update to %s", req.Status)
	}
	if err := h.adminStore.TransitionStatus(c.Request.Context(), id, req.Status, operatorID, req.Remark); err != nil {
		if err == repository.ErrInvalidOrderStatusTransition {
			response.Error(c, http.StatusConflict, errcode.ErrConflict, err.Error())
			return
		}
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

// UpdateStatus 兼容旧调用，统一转发到状态机入口。
func (h *BookingHandler) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	if h.adminStore == nil {
		return fmt.Errorf("booking store unavailable")
	}
	return h.adminStore.TransitionStatus(ctx, orderID, status, 0, "handler forward")
}

func parseOperatorID(c *gin.Context) int64 {
	if value, ok := c.Get(middleware.ContextKeyStaffID); ok {
		id, err := strconv.ParseInt(fmt.Sprint(value), 10, 64)
		if err == nil && id > 0 {
			return id
		}
	}
	return 0
}

// AdminExport 管理后台导出订单 CSV。
func (h *BookingHandler) AdminExport(c *gin.Context) {
	if h.exportService == nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "export service unavailable")
		return
	}

	filter := service.OrderFilter{
		Status:    c.Query("status"),
		BookingNo: c.Query("booking_no"),
	}
	if voyageID, err := strconv.ParseInt(c.Query("voyage_id"), 10, 64); err == nil {
		filter.VoyageID = voyageID
	}

	ctx := service.WithOrderExportPermission(c.Request.Context(), true)
	data, err := h.exportService.ExportToExcel(ctx, filter)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	filename := fmt.Sprintf("orders_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, "text/csv", data)
}
